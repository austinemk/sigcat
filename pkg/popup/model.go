package popup

import (
	_ "embed"

	tea "charm.land/bubbletea/v2"
	"github.com/austinemk/sigcat/pkg/helpers"
)

//go:embed cat.txt
var catASCII string

type PopupModel struct {
	Task          helpers.BreakTask
	DaemonRunning bool
}

// InitialPopupModel looks up the relevant active task parameters or supplies fallback presets
func InitialPopupModel(id string) PopupModel {
	tasks, _ := helpers.LoadTasks()
	var targeted helpers.BreakTask

	for _, t := range tasks {
		if t.ID == id {
			targeted = t
			break
		}
	}

	if targeted.ID == "" {
		targeted = helpers.BreakTask{
			Title:   "Take a Break!",
			Message: "Time to stretch and look away.",
		}
	}

	return PopupModel{
		Task:          targeted,
		DaemonRunning: helpers.IsDaemonRunning(),
	}
}

func (m PopupModel) Init() tea.Cmd {
	return nil
}
