package tui

import (
	"fmt"
	"os"

	"github.com/ChmaraX/notidb/internal"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type saveModel struct {
	dbId    string
	entry   internal.DatabaseEntry
	saving  bool
	res     res
	spinner spinner.Model
}

type res struct {
	url string
	err error
}

func (m saveModel) saveEntry() tea.Msg {
	page, err := internal.AddDatabaseEntry(m.dbId, m.entry)
	if err != nil {
		return res{url: "", err: fmt.Errorf("failed to add database entry: %w", err)}
	}
	return res{url: page.URL, err: nil}
}

func (m saveModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.saveEntry)
}

func (m saveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case res:
		m.saving = false
		m.res = msg
		return m, tea.Quit
	}

	return m, cmd
}

func (m saveModel) View() string {
	switch {
	case m.res.err != nil:
		return quitTextStyle.Render("Error: " + m.res.err.Error())
	case !m.saving:
		return quitTextStyle.Render("Saved: " + m.res.url)
	default:
		return m.spinner.View() + "Saving..."
	}
}

func NewSaveModel(dbId string, entry internal.DatabaseEntry) *saveModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return &saveModel{dbId: dbId, entry: entry, spinner: s, saving: true}
}

func InitSave(dbId string, entry internal.DatabaseEntry) {
	m := NewSaveModel(dbId, entry)
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
