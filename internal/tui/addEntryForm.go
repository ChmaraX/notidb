package tui

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/ChmaraX/notidb/internal"
	"github.com/ChmaraX/notidb/internal/utils"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jomei/notionapi"
)

func InitForm(schema notionapi.PropertyConfigs) {
	p := tea.NewProgram(initialModel(schema))

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

func cvvValidator(s string) error {
	// The CVV should be a number of 3 digits
	// Since the input will already ensure that the CVV is a string of length 3,
	// All we need to do is check that it is a number
	_, err := strconv.ParseInt(s, 10, 64)
	return err
}

func initialModel(schema notionapi.PropertyConfigs) model {
	var inputs []textinput.Model = make([]textinput.Model, len(schema))
	supportedPropTypes := internal.GetSupportedPagePropTypes()

	// TODO: append Title and Content to the schema
	// TODO: body = textfield/editor
	// TOOD: shift+enter submit?

	// iterate over schema and create inputs, but only for supported types
	idx := 0
	for title, prop := range schema {
		if utils.Contains(supportedPropTypes, string(prop.GetType())) {
			switch prop.GetType() {
			case
				notionapi.PropertyConfigTypeRichText,
				notionapi.PropertyConfigTypeNumber, // TOOD: number validator
				notionapi.PropertyConfigTypeSelect,
				notionapi.PropertyConfigTypeMultiSelect,
				notionapi.PropertyConfigTypeDate, // TODO: date validator
				notionapi.PropertyConfigTypeCheckbox,
				notionapi.PropertyConfigTypeURL, // TODO: url validator
				notionapi.PropertyConfigTypeEmail,
				notionapi.PropertyConfigTypePhoneNumber: // TODO: phone number validator
				inputs[idx] = textinput.New()
				inputs[idx].Placeholder = title
			default:
				fmt.Printf("unsupported property type: %s", prop.GetType())
			}
			idx++
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
