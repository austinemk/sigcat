package dashboard

import (
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/austinemk/sigcat/pkg/helpers"
)

type sessionState int

const (
	viewTasks sessionState = iota
	createTask
)

// ==========================================
// Core Domain Structural Types
// ==========================================

type dashboardModel struct {
	state         sessionState
	taskList      list.Model
	inputIndex    int
	inputs        []textinput.Model
	errMessage    string
	daemonRunning bool
}

// ==========================================
// Localized Task List Delegate
// ==========================================

type taskDelegate struct{}

func (d taskDelegate) Height() int                               { return 2 }
func (d taskDelegate) Spacing() int                              { return 1 }
func (d taskDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

// ==========================================
// Initialization Loop & Update Loop Router
// ==========================================

func InitialDashboardModel() dashboardModel {
	t, _ := helpers.LoadTasks()

	listItems := make([]list.Item, len(t))
	for i, task := range t {
		listItems[i] = task
	}

	taskGrid := list.New(listItems, taskDelegate{}, 0, 0)
	taskGrid.SetShowTitle(false)
	taskGrid.SetShowStatusBar(false)
	taskGrid.SetShowHelp(false)

	inputs := make([]textinput.Model, 4)
	inputs[0] = textinput.New()
	inputs[0].Placeholder = "Stretch Break"
	inputs[0].Focus()

	inputs[1] = textinput.New()
	inputs[1].Placeholder = "Look away for 20 seconds!"

	inputs[2] = textinput.New()
	inputs[2].Placeholder = "20"

	inputs[3] = textinput.New()
	inputs[3].Placeholder = "y/n"

	return dashboardModel{
		state:         viewTasks,
		taskList:      taskGrid,
		inputs:        inputs,
		daemonRunning: helpers.IsDaemonRunning(),
	}
}

func (m dashboardModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m dashboardModel) getTasks() []helpers.BreakTask {
	items := m.taskList.Items()
	tasks := make([]helpers.BreakTask, len(items))
	for i, item := range items {
		tasks[i] = item.(helpers.BreakTask)
	}
	return tasks
}

func (m dashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.daemonRunning = helpers.IsDaemonRunning()
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if updatedModel, handled, globalCmd := m.handleGlobalKeys(msg); handled {
			return updatedModel, globalCmd
		}

		if m.state == viewTasks {
			return m.handleViewTasksKeys(msg)
		} else if m.state == createTask {
			return m.handleCreateTaskKeys(msg)
		}
	}

	if m.state == createTask {
		m.inputs[m.inputIndex], cmd = m.inputs[m.inputIndex].Update(msg)
		return m, cmd
	}

	return m, nil
}
