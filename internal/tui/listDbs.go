package tui

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ChmaraX/notidb/internal"
	"github.com/ChmaraX/notidb/internal/settings"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item struct {
	id, title string
	def       bool
}

var showIds bool

func (i item) FilterValue() string { return string(i.title) }

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

type model struct {
	list        list.Model
	choice      string
	quitting    bool
	defaultDbId string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				m.choice = string(i.title)
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

func (m model) View() string {
	if m.choice != "" {
		return quitTextStyle.Render(fmt.Sprintf("Database '%s' successfuly set as default.", m.choice))
	}
	if m.quitting {
		if m.defaultDbId == settings.NoDefaultDatabaseId {
			return quitTextStyle.Render("No default database set.")
		}
		return quitTextStyle.Render("No changes made.")
	}
	return "\n" + m.list.View()
}

func GetDbsListTUI(dbs []internal.NotionDb, defaultDbId string) {
	items := make([]list.Item, len(dbs))
	for i, db := range dbs {
		items[i] = item{title: db.Title, id: db.Id, def: db.Id == defaultDbId}
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Choose default database for operations:"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{list: l, defaultDbId: defaultDbId}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
