package tui

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ChmaraX/notidb/internal/notion"
	"github.com/ChmaraX/notidb/internal/utils"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jomei/notionapi"
)

func InitForm(schema notionapi.PropertyConfigs) notion.DatabaseEntry {
	props := filterSupportedProps(schema)
	model, err := tea.NewProgram(initialModel(props)).Run()

	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	return model.(formModel).entry
}

const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
	darkRed  = lipgloss.Color("#E05252")
)

var (
	inputStyle = lipgloss.NewStyle().Foreground(hotPink).MarginLeft(2)
	errorStyle = lipgloss.NewStyle().Foreground(darkRed).MarginLeft(2)
)

var validators = map[notionapi.PropertyType]textinput.ValidateFunc{
	notionapi.PropertyTypeNumber:   numberValidator,
	notionapi.PropertyTypeCheckbox: checkboxValidator,
}

var placeholders = map[notionapi.PropertyType]string{
	notionapi.PropertyTypeRichText:    "Enter text",
	notionapi.PropertyTypeSelect:      "Enter value",
	notionapi.PropertyTypeMultiSelect: "Enter comma separated values",
	notionapi.PropertyTypeDate:        "31/12/1990",
	notionapi.PropertyTypeCheckbox:    "y/n",
	notionapi.PropertyTypeNumber:      "123",
	notionapi.PropertyTypeEmail:       "example@email.com",
	notionapi.PropertyTypePhoneNumber: "+48 123 456 789",
}

type formModel struct {
	entry        notion.DatabaseEntry
	props        []PropInput
	block        BlockInput
	focusedProp  int
	focusOnProps bool
	err          error
	help         help.Model
	keymap       keymap
}

type PropInput struct {
	propType notionapi.PropertyType
	model    textinput.Model
	title    string
}

type BlockInput struct {
	propType notionapi.Block
	model    textarea.Model
}

type keymap struct {
	save key.Binding
	next key.Binding
	prev key.Binding
	quit key.Binding
}

func (m formModel) toDatabaseEntry() (notion.DatabaseEntry, error) {
	entry := notion.DatabaseEntry{
		Props:  make(notionapi.Properties),
		Blocks: make([]notionapi.Block, 0),
	}

	for _, prop := range m.props {
		propTitle := prop.title
		propValue := prop.model.Value()

		if propValue == "" {
			continue
		}

		switch prop.propType {
		case notionapi.PropertyTypeTitle:
			entry.Props[propTitle] = notion.CreateTitleProperty(propValue)
		case notionapi.PropertyTypeRichText:
			entry.Props[propTitle] = notion.CreateRichTextProperty(propValue)
		case notionapi.PropertyTypeSelect:
			entry.Props[propTitle] = notion.CreateSelectProperty(propValue)
		case notionapi.PropertyTypeMultiSelect:
			v := strings.Split(propValue, ",")
			entry.Props[propTitle] = notion.CreateMultiSelectProperty(v)
		case notionapi.PropertyTypeDate:
			date, err := notion.CreateDateProperty(propValue)
			if err != nil {
				return notion.DatabaseEntry{}, fmt.Errorf("failed to parse date: %w", err)
			}
			entry.Props[propTitle] = date
		case notionapi.PropertyTypeCheckbox:
			v, _ := utils.ParseBool(propValue)
			entry.Props[propTitle] = notion.CreateCheckboxProperty(v)
		case notionapi.PropertyTypeNumber:
			v, _ := strconv.ParseFloat(propValue, 64)
			entry.Props[propTitle] = notion.CreateNumberProperty(v)
		case notionapi.PropertyTypeEmail:
			entry.Props[propTitle] = notion.CreateEmailProperty(propValue)
		case notionapi.PropertyTypePhoneNumber:
			entry.Props[propTitle] = notion.CreatePhoneNumberProperty(propValue)
		default:
			fmt.Printf("unsupported property type: %s", prop.propType)
		}
	}

	entry.Blocks = append(entry.Blocks, notion.CreateContentBlock(m.block.model.Value()))

	return entry, nil
}

// filter props supported by the TUI form
func filterSupportedProps(schema notionapi.PropertyConfigs) map[string]notionapi.PropertyType {
	supportedPropTypes := notion.GetSupportedPropTypes()

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

func numberValidator(s string) error {
	_, err := strconv.ParseFloat(s, 64)
	if err != nil && s != "" {
		return fmt.Errorf("must be number")
	}
	return nil
}

func checkboxValidator(s string) error {
	_, err := utils.ParseBool(s)
	if err != nil && s != "" {
		return fmt.Errorf("must be y/n")
	}
	return nil
}

func createPropInput(title string, propType notionapi.PropertyType) PropInput {
	ti := textinput.New()
	ti.Placeholder = placeholders[propType]
	ti.Validate = validators[propType]

	return PropInput{
		propType: propType,
		model:    ti,
		title:    title,
	}
}

func createBlockInput() BlockInput {
	ta := textarea.New()
	ta.Placeholder = "Start writing..."
	ta.SetWidth(50)
	ta.SetHeight(10)

	return BlockInput{
		propType: notionapi.ParagraphBlock{},
		model:    ta,
	}
}

func initialModel(props map[string]notionapi.PropertyType) formModel {
	// create map of prop inputs
	propInputs := make([]PropInput, len(props))

	titleIdx := 0 // title is always first
	idx := 1
	for title, propType := range props {

		switch propType {
		case notionapi.PropertyTypeTitle:
			pi := createPropInput(title, propType)
			pi.model.Focus()
			propInputs[titleIdx] = pi
		case
			notionapi.PropertyTypeRichText,
			notionapi.PropertyTypeSelect,
			notionapi.PropertyTypeMultiSelect,
			notionapi.PropertyTypeDate,
			notionapi.PropertyTypeCheckbox,
			notionapi.PropertyTypeNumber,
			notionapi.PropertyTypeEmail,
			notionapi.PropertyTypePhoneNumber:
			propInputs[idx] = createPropInput(title, propType)
			idx++
		default:
			fmt.Printf("unsupported property type: %s", propType)
		}
	}

	// help styles
	help := help.New()
	help.Styles.ShortKey = lipgloss.NewStyle().Foreground(darkGray)

	return formModel{
		props:        propInputs,
		block:        createBlockInput(),
		focusedProp:  0,
		focusOnProps: true,
		err:          nil,
		keymap:       getHelpKeyMap(),
		help:         help,
	}
}

func getHelpKeyMap() keymap {
	return keymap{
		save: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("<ctrl+s>", "save"),
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

func (m formModel) helpView() string {
	return m.help.ShortHelpView([]key.Binding{
		m.keymap.save,
		m.keymap.next,
		m.keymap.prev,
		m.keymap.quit,
	})
}

func (m formModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m formModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.props)+1) // +1 for block input

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlS:
			entry, err := m.toDatabaseEntry()
			if err != nil {
				m.err = err
				return m, nil
			}
			m.entry = entry
			return m, tea.Quit
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		}
	}

	// Update each element and collect commands
	for i := range m.props {
		m.props[i].model, cmds[i] = m.props[i].model.Update(msg)
	}
	m.block.model, cmds[len(m.props)] = m.block.model.Update(msg)

	return m, tea.Batch(cmds...)
}

func getElemErrMsg(input textinput.Model) string {
	if input.Err != nil {
		return errorStyle.Render(input.Err.Error())
	}
	return ""
}

func (m formModel) View() string {
	var inputsView strings.Builder

	for _, value := range m.props {
		input := value.model
		inputsView.WriteString(fmt.Sprintf("%s%s%s\n", inputStyle.Width(15).Render(value.title), input.View(), getElemErrMsg(input)))
	}

	inputsView.WriteString(fmt.Sprintf("\n%s\n%s\n", inputStyle.Width(30).Render("Content"), m.block.model.View()))

	if m.err != nil {
		inputsView.WriteString(fmt.Sprintf("\n%s\n", errorStyle.Render(m.err.Error())))
	}

	return fmt.Sprintf("\n%s\n%s\n\n", inputsView.String(), m.helpView())
}

func (m *formModel) nextInput() {
	m.blurCurrentElement()

	if m.focusOnProps {
		m.focusedProp++
		if m.focusedProp >= len(m.props) {
			m.focusOnProps = false
			m.focusedProp = 0
		}
	} else {
		m.focusOnProps = true
	}

	m.focusCurrentElement()
}

func (m *formModel) prevInput() {
	m.blurCurrentElement()

	if m.focusOnProps {
		if m.focusedProp == 0 {
			m.focusOnProps = false
		} else {
			m.focusedProp--
		}
	} else {
		m.focusOnProps = true
		m.focusedProp = len(m.props) - 1
	}

	m.focusCurrentElement()
}

func (m *formModel) blurCurrentElement() {
	if m.focusOnProps {
		m.props[m.focusedProp].model.Blur()
	} else {
		m.block.model.Blur()
	}
}

func (m *formModel) focusCurrentElement() {
	if m.focusOnProps {
		m.props[m.focusedProp].model.Focus()
	} else {
		m.block.model.Focus()
	}
}
