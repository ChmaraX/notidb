package tui

import (
	"fmt"
	"os"

	"github.com/ChmaraX/notidb/internal"
	"github.com/ChmaraX/notidb/internal/settings"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func loadDatabases() tea.Msg {
	databases, err := internal.GetAllNotionDbs()
	if err != nil {
		return dbs{list: nil, err: fmt.Errorf("error loading databases: %v", err), loaded: true}
	}
	if len(databases) == 0 {
		return dbs{list: nil, err: fmt.Errorf("no databases found in your workspace or the access is not granted"), loaded: true}
	}
	return dbs{list: databases, err: nil, loaded: true}
}

func loadDefaultDatabase() tea.Msg {
	defaultDbId, err := settings.GetDefaultDatabase()
	if err != nil {
		return defaultDb{id: "", err: fmt.Errorf("error loading default database: %v", err), loaded: true}
	}
	return defaultDb{id: defaultDbId, err: nil, loaded: true}
}

type dbs struct {
	list   []internal.NotionDb
	err    error
	loaded bool
}

type defaultDb struct {
	id     string
	err    error
	loaded bool
}

type loadingModel struct {
	spinner   spinner.Model
	dbs       dbs
	defaultDb defaultDb
	err       error
}

func newLoadingModel() loadingModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return loadingModel{spinner: s}
}

func checkAllLoaded(m loadingModel) (tea.Model, tea.Cmd) {
	dbs, defaultDbId := m.dbs.list, m.defaultDb.id

	if !m.dbs.loaded || !m.defaultDb.loaded {
		return m, nil
	}

	if !dbExists(dbs, defaultDbId) && defaultDbId != settings.NoDefaultDatabaseId {
		m.err = fmt.Errorf("database which is set as default (%s) was not found in your workspace or the access is not granted", defaultDbId)
		return m, tea.Quit
	}

	dbListModel := InitDbListModel(dbs, defaultDbId)
	return dbListModel.Update(nil)
}

func dbExists(dbs []internal.NotionDb, dbId string) bool {
	for _, db := range dbs {
		if db.Id == dbId {
			return true
		}
	}
	return false
}

func (m loadingModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, loadDatabases, loadDefaultDatabase)
}

func (m loadingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case dbs:
		m.dbs = msg
		return checkAllLoaded(m)
	case defaultDb:
		m.defaultDb = msg
		return checkAllLoaded(m)
	}
	return m, nil
}

func (m loadingModel) View() string {
	switch {
	case m.dbs.err != nil:
		return quitTextStyle.Render("Error: " + m.dbs.err.Error())
	case m.defaultDb.err != nil:
		return quitTextStyle.Render("Error: " + m.defaultDb.err.Error())
	case m.err != nil:
		return quitTextStyle.Render("Error: " + m.err.Error())
	case !m.dbs.loaded:
		return m.spinner.View() + " Calling Notion API - loading databases..."
	case !m.defaultDb.loaded:
		return m.spinner.View() + " Reading settings..."
	default:
		return quitTextStyle.Render("Databases loaded successfully.")
	}
}

func LoadDbs() {
	m := newLoadingModel()
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
