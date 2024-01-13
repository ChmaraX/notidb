package tui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type LoadingModel struct {
	spinner    spinner.Model
	asyncFuncs []func() tea.Msg
	NumFuncs   int
	err        error
	Responses  []Response
}

type LoadingFunc func() Response

func (f LoadingFunc) wrapAsMsg() func() tea.Msg {
	return func() tea.Msg {
		result := f()
		return Response{
			Id:   result.Id,
			Data: result.Data,
			Err:  result.Err,
		}
	}
}

func mapFuncsToMsgs(funcs []LoadingFunc) []func() tea.Msg {
	msgs := make([]func() tea.Msg, len(funcs))
	for i, f := range funcs {
		msgs[i] = f.wrapAsMsg()
	}
	return msgs
}

type Response struct {
	Id   string
	Data interface{}
	Err  error
}

func newLoadingModel1(funcs ...LoadingFunc) LoadingModel {
	s := spinner.New()
	s.Spinner = spinner.Dot

	return LoadingModel{spinner: s, asyncFuncs: mapFuncsToMsgs(funcs), NumFuncs: len(funcs)}
}

func (m LoadingModel) GetResponse(id string) Response {
	for _, item := range m.Responses {
		if item.Id == id {
			return item
		}
	}
	return Response{}
}

func (m LoadingModel) Init() tea.Cmd {
	cmds := make([]tea.Cmd, len(m.asyncFuncs))
	for i, f := range m.asyncFuncs {
		cmds[i] = tea.Cmd(f)
	}
	return tea.Batch(append(cmds, m.spinner.Tick)...)
}

func (m LoadingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case Response:
		m.Responses = append(m.Responses, msg)
		m.err = msg.Err

		if len(m.Responses) == m.NumFuncs {
			return m, tea.Quit
		}

		return m, nil
	}
	return m, nil
}

func (m LoadingModel) View() string {

	if m.err != nil {
		return fmt.Sprintf("Error loading: %v", m.err)
	}

	if len(m.Responses) == m.NumFuncs {
		return "Done loading"
	}

	return m.spinner.View() + " Loading..."
}

func NewLoadingModel(funcs ...LoadingFunc) LoadingModel {
	m := newLoadingModel1(funcs...)
	model, err := tea.NewProgram(m).Run()

	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	return model.(LoadingModel)
}
