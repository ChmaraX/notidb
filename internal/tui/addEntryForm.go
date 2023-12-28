package tui

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/ChmaraX/notidb/internal"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jomei/notionapi"
)

func InitForm(dbId string) {
	entryForm := createEntryInputForm(dbId)
	p := tea.NewProgram(initialModel(entryForm))

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg error
)

const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
)

var (
	inputStyle = lipgloss.NewStyle().Foreground(hotPink).MarginLeft(2)
)

type model struct {
	elements   []interface{} // common array for all inputs
	focused    int
	err        error
	bottomText string
	help       help.Model
	keymap     keymap
}

type keymap struct {
	submit key.Binding
	next   key.Binding
	prev   key.Binding
	quit   key.Binding
}

type EntryInputForm struct {
	Props   map[string]notionapi.PropertyType
	Content string
}

func createEntryInputForm(dbId string) EntryInputForm {
	schema, err := internal.GetDatabaseSchema(dbId)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	props := filterSupportedProps(schema)

	return EntryInputForm{
		Props:   props,
		Content: "",
	}
}

// filter props supported by the TUI form
func filterSupportedProps(schema notionapi.PropertyConfigs) map[string]notionapi.PropertyType {
	supportedPropTypes := internal.GetSupportedPropTypes()

	// Convert slice to map for faster lookup
	supportedPropTypesMap := make(map[string]bool)
	for _, propType := range supportedPropTypes {
		supportedPropTypesMap[string(propType)] = true
	}

	props := make(map[string]notionapi.PropertyType)
	for key, prop := range schema {
		// Use map lookup instead of slice contains
		if _, ok := supportedPropTypesMap[string(prop.GetType())]; ok {
			props[key] = notionapi.PropertyType(prop.GetType())
		}
	}

	return props
}

func expValidator(s string) error {
	// The 3 character should be a slash (/)
	// The rest should be numbers
	e := strings.ReplaceAll(s, "/", "")
	_, err := strconv.ParseInt(e, 10, 64)
	if err != nil {
		return fmt.Errorf("EXP is invalid")
	}

	// There should be only one slash and it should be in the 2nd index (3rd character)
	if len(s) >= 3 && (strings.Index(s, "/") != 2 || strings.LastIndex(s, "/") != 2) {
		return fmt.Errorf("EXP is invalid")
	}

	return nil
}

func initialModel(entryForm EntryInputForm) model {
	elements := make([]interface{}, len(entryForm.Props)+1)

	titleIdx := 0 // always first
	idx := 1
	for title, prop := range entryForm.Props {
		switch prop {
		case notionapi.PropertyTypeTitle:
			ti := textinput.New()
			ti.Placeholder = title
			ti.Focus()

			elements[titleIdx] = ti
		case
			notionapi.PropertyTypeRichText,
			notionapi.PropertyTypeNumber, // TOOD: number validator
			notionapi.PropertyTypeSelect,
			notionapi.PropertyTypeMultiSelect,
			notionapi.PropertyTypeDate, // TODO: date validator
			notionapi.PropertyTypeCheckbox,
			notionapi.PropertyTypeEmail,
			notionapi.PropertyTypePhoneNumber: // TODO: phone number validator

			ti := textinput.New()
			ti.Placeholder = title

			elements[idx] = ti
			idx++
		default:
			fmt.Printf("unsupported property type: %s", prop)
		}
	}

	ta := textarea.New()
	ta.Placeholder = "Start writing..."
	ta.SetWidth(50)
	ta.SetHeight(10)

	elements[len(elements)-1] = ta

	return model{
		elements: elements,
		focused:  0,
		err:      nil,
		keymap:   getHelpKeyMap(),
		help:     help.New(),
	}
}

func getHelpKeyMap() keymap {
	return keymap{
		submit: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("<enter>", "submit"),
		),
		next: key.NewBinding(
			key.WithKeys("tab", "ctrl+n"),
			key.WithHelp("<tab>", "next"),
		),
		prev: key.NewBinding(
			key.WithKeys("shift+tab", "ctrl+p"),
			key.WithHelp("<shift+tab>", "previous"),
		),
		quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("<ctrl+c>", "quit"),
		),
	}
}

func (m model) helpView() string {
	return m.help.ShortHelpView([]key.Binding{
		m.keymap.submit,
		m.keymap.next,
		m.keymap.prev,
		m.keymap.quit,
	})
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.elements))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return m, tea.Quit
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	// Update each element and collect commands
	for i := range m.elements {
		switch elem := m.elements[i].(type) {
		case textinput.Model:
			m.elements[i], cmds[i] = elem.Update(msg)
		case textarea.Model:
			m.elements[i], cmds[i] = elem.Update(msg)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var inputsView strings.Builder

	for i, elem := range m.elements {
		switch elem := elem.(type) {
		case textinput.Model:
			inputsView.WriteString(fmt.Sprintf("%s\n%s\n", inputStyle.Width(30).Render(elem.Placeholder), elem.View()))
		case textarea.Model:
			if i == len(m.elements)-1 { // TextArea is the last element
				inputsView.WriteString(fmt.Sprintf("\n%s\n%s\n", inputStyle.Width(30).Render("Content"), elem.View()))
			}
		}
	}

	return fmt.Sprintf(
		"\n%s\n%s\n",
		inputsView.String(),
		m.helpView(),
	) + "\n"
}

func (m *model) nextInput() {
	m.blurCurrentElement()
	m.focused = (m.focused + 1) % len(m.elements)
	m.focusCurrentElement()
}

func (m *model) prevInput() {
	m.blurCurrentElement()
	m.focused = (m.focused - 1 + len(m.elements)) % len(m.elements)
	m.focusCurrentElement()
}

func (m *model) blurCurrentElement() {
	if elem, ok := m.elements[m.focused].(textinput.Model); ok {
		elem.Blur()
		m.elements[m.focused] = elem
	} else if elem, ok := m.elements[m.focused].(textarea.Model); ok {
		elem.Blur()
		m.elements[m.focused] = elem
	}
}

func (m *model) focusCurrentElement() {
	if elem, ok := m.elements[m.focused].(textinput.Model); ok {
		elem.Focus()
		m.elements[m.focused] = elem
	} else if elem, ok := m.elements[m.focused].(textarea.Model); ok {
		elem.Focus()
		m.elements[m.focused] = elem
	}
}
