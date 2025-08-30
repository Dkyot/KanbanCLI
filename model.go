package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	list  list.Model
	err   error
	ready bool // flag to initialize once
}

func New() *Model {
	return &Model{}
}

func (m *Model) initList(width, height int) {
	delegate := list.NewDefaultDelegate()
	items := []list.Item{
		Task{status: todo, title: "buy milk", description: "strawberry milk"},
		Task{status: todo, title: "eat sushi", description: "negitoro roll"},
		Task{status: todo, title: "fold laundry", description: "or wear wrinkly t-shirt"},
	}

	l := list.New(items, delegate, width, height)
	l.Title = "To Do"
	m.list = l
	m.ready = true
}

// region model interface implementation

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.ready {
			m.initList(msg.Width, msg.Height)
		} else {
			m.list.SetSize(msg.Width, msg.Height)
		}
	}

	if m.ready {
		m.list, cmd = m.list.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	if !m.ready {
		return "loading..."
	}
	return m.list.View()
}

// endregion
