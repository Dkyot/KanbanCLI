package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// some settings to render
const divisor = 4

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
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width/divisor, height)
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
			m.lists[m.focused].SetSize(msg.Width, msg.Height)
		}
	}

	var cmd tea.Cmd
	if m.ready {
		m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	if !m.ready {
		return "loading..."
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.lists[todo].View(),
		m.lists[inProgress].View(),
		m.lists[done].View(),
	)
}

// endregion
