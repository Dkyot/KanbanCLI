package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// some settings to render
const divisor = 3
const rightOffset = 2
const verticalOffset = 4

var (
	columnStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.HiddenBorder())
		//Background(lipgloss.Color("11"))
	focusedStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))
		//Background(lipgloss.Color("10"))
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

type Model struct {
	lists   []list.Model
	focused status
	err     error
	ready   bool // flag to initialize once
}

func New() *Model {
	return &Model{}
}

func (m *Model) initLists(width, height int) {
	colWidth := (width / divisor) - rightOffset
	columnStyle = columnStyle.Width(colWidth)
	focusedStyle = focusedStyle.Width(colWidth)

	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), colWidth, height-verticalOffset)
	defaultList.SetShowHelp(false) // remove void space from right side
	m.lists = []list.Model{defaultList, defaultList, defaultList}

	// init To Do
	m.lists[todo].Title = "To Do"
	m.lists[todo].SetItems([]list.Item{
		Task{status: todo, title: "buy milk", description: "strawberry milk"},
		Task{status: todo, title: "eat sushi", description: "negitoro roll"},
		Task{status: todo, title: "fold laundry", description: "or wear wrinkly t-shirt"},
	})

	// init In Progress
	m.lists[inProgress].Title = "In Progress"
	m.lists[inProgress].SetItems([]list.Item{
		Task{status: inProgress, title: "create golang project", description: "learn about bubble tea"},
	})

	// init Done
	m.lists[done].Title = "Done"
	m.lists[done].SetItems([]list.Item{
		Task{status: done, title: "learn go", description: "learn golang syntax"},
	})

	m.ready = true
}

// region switch active columns

func (m *Model) Next() {
	if m.focused == done {
		m.focused = todo
	} else {
		m.focused++
	}
}

func (m *Model) Prev() {
	if m.focused == todo {
		m.focused = done
	} else {
		m.focused--
	}
}

// endregion

// region model interface implementation

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.ready {
			m.initLists(msg.Width, msg.Height)
		} else {
			colWidth := (msg.Width / divisor) - rightOffset
			columnStyle = columnStyle.Width(colWidth)
			focusedStyle = focusedStyle.Width(colWidth)

			for i := range m.lists {
				m.lists[i].SetSize(colWidth, msg.Height-verticalOffset)
			}
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			m.Prev()
		case "right", "l":
			m.Next()
		}
	}

	if m.ready {
		cmds := make([]tea.Cmd, len(m.lists))
		for i := range m.lists {
			var c tea.Cmd
			m.lists[i], c = m.lists[i].Update(msg)
			cmds[i] = c
		}
		return m, tea.Batch(cmds...)
	}

	return m, nil
}

func (m Model) View() string {
	if !m.ready {
		return "loading..."
	}

	todoView := m.lists[todo].View()
	inProgressView := m.lists[inProgress].View()
	doneiew := m.lists[done].View()

	var renderStr string
	switch m.focused {
	case todo:
		renderStr = lipgloss.JoinHorizontal(
			lipgloss.Left,
			focusedStyle.Render(todoView),
			columnStyle.Render(inProgressView),
			columnStyle.Render(doneiew),
		)
	case inProgress:
		renderStr = lipgloss.JoinHorizontal(
			lipgloss.Left,
			columnStyle.Render(todoView),
			focusedStyle.Render(inProgressView),
			columnStyle.Render(doneiew),
		)
	case done:
		renderStr = lipgloss.JoinHorizontal(
			lipgloss.Left,
			columnStyle.Render(todoView),
			columnStyle.Render(inProgressView),
			focusedStyle.Render(doneiew),
		)
	default:
		renderStr = lipgloss.JoinHorizontal(
			lipgloss.Left,
			columnStyle.Render(todoView),
			columnStyle.Render(inProgressView),
			columnStyle.Render(doneiew),
		)
	}

	return renderStr
}

// endregion
