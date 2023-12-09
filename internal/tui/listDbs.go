package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/ChmaraX/notidb/internal"
	"github.com/ChmaraX/notidb/internal/settings"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14
const listWidth = 60

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	checkMark         = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render("✓")
)

type item struct {
	id, title string
	def       bool
}

var showIds = false

func (i item) FilterValue() string { return i.title }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%d. %s", index+1, i.title))

	if i.def {
		builder.WriteString(" [default]")
	}

	if showIds {
		builder.WriteString(fmt.Sprintf(" (%s)", i.id))
	}

	var renderFn func(...string) string
	if index == m.Index() {
		renderFn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	} else {
		renderFn = itemStyle.Render
	}

	fmt.Fprint(w, renderFn(builder.String()))
}

type DbListModel struct {
	list        list.Model
	choice      string
	quitting    bool
	defaultDbId string
}

func (m DbListModel) Init() tea.Cmd {
	return nil
}

func (m DbListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = i.title
				settings.SetDefaultDatabase(i.id)
			}
			return m, tea.Quit

		case "c":
			showIds = !showIds
			return m, nil
		}

	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m DbListModel) View() string {
	if m.choice != "" {
		highlightStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("170"))
		return quitTextStyle.Render(fmt.Sprintf("%s Default database successfully set to: %s", checkMark, highlightStyle.Render(m.choice)))
	}
	if m.quitting {
		if m.defaultDbId == settings.NoDefaultDatabaseId {
			return quitTextStyle.Render("No default database set.")
		}
		return quitTextStyle.Render("No changes made.")
	}
	return "\n" + m.list.View()
}

func InitDbListModel(dbs []internal.NotionDb, defaultDbId string) *DbListModel {
	items := make([]list.Item, len(dbs))
	for i, db := range dbs {
		items[i] = item{title: db.Title, id: db.Id, def: db.Id == defaultDbId}
	}

	l := initListModel(items)
	m := DbListModel{list: l, defaultDbId: defaultDbId}

	return &m
}

func initListModel(items []list.Item) list.Model {
	l := list.New(items, itemDelegate{}, listWidth, listHeight)
	l.Title = "Choose default database for operations:"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return l
}
