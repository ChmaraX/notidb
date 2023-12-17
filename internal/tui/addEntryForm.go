package tui

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/ChmaraX/notidb/internal"
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
	inputStyle      = lipgloss.NewStyle().Foreground(hotPink).MarginLeft(2)
	bottomTextStyle = lipgloss.NewStyle().Foreground(darkGray)
)

type model struct {
	inputs     []textinput.Model
	focused    int
	err        error
	bottomText string
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
	var inputs []textinput.Model = make([]textinput.Model, len(entryForm.Props))

	// TODO: body = textfield/editor
	// TOOD: shift+enter submit?

	idx := 1
	for title, prop := range entryForm.Props {
		switch prop {
		case
			notionapi.PropertyTypeTitle:
			inputs[0] = textinput.New()
			inputs[0].Placeholder = title
			inputs[0].Focus()
		case
			notionapi.PropertyTypeRichText,
			notionapi.PropertyTypeNumber, // TOOD: number validator
			notionapi.PropertyTypeSelect,
			notionapi.PropertyTypeMultiSelect,
			notionapi.PropertyTypeDate, // TODO: date validator
			notionapi.PropertyTypeCheckbox,
			notionapi.PropertyTypeEmail,
			notionapi.PropertyTypePhoneNumber: // TODO: phone number validator
			inputs[idx] = textinput.New()
			inputs[idx].Placeholder = title
			idx++
		default:
			fmt.Printf("unsupported property type: %s", prop)
		}
	}

	return model{
		inputs:     inputs,
		focused:    0,
		err:        nil,
		bottomText: "Continue ->",
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.focused == len(m.inputs)-1 {
				// print the values
				for i := range m.inputs {
					// TODO: create a body and send in POST request
					fmt.Printf("%s: %s\n", m.inputs[i].Placeholder, m.inputs[i].Value())
				}
				return m, tea.Quit
			}
			m.nextInput()
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
			m.updateBottomText()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
			m.updateBottomText()
		}
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.focused].Focus()

	case errMsg:
		m.err = msg
		return m, nil
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var inputsView string
	for i := range m.inputs {
		inputsView += fmt.Sprintf("%s\n%s\n", inputStyle.Width(30).Render(m.inputs[i].Placeholder), m.inputs[i].View())
	}

	return fmt.Sprintf(
		` %s

%s

 %s
`,
		bottomTextStyle.Width(30).Bold(true).Render("Add new entry to database:"),

		inputsView,

		bottomTextStyle.Render(m.bottomText),
	) + "\n"
}

func (m *model) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

func (m *model) prevInput() {
	m.focused--
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}

func (m *model) updateBottomText() {
	if m.focused == len(m.inputs)-1 {
		m.bottomText = "Press Enter to Submit"
	} else {
		m.bottomText = "Continue ->"
	}
}
