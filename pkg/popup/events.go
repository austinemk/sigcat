package popup

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/austinemk/sigcat/pkg/helpers"
)

func (m PopupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "s", "q", "escape":
			return m, tea.Quit

		case "r":
			tasks, err := helpers.LoadTasks()
			if err == nil {
				for i, t := range tasks {
					if t.ID == m.Task.ID {
						if t.AutoRepeat {
							// Turn the repeat sequence loop off completely if requested
							tasks[i].AutoRepeat = false
						} else {
							// Shift time constraints onward to postpone safely
							tasks[i].IsActive = true
							tasks[i].NextRun = time.Now().Add(time.Duration(t.DurationMin) * time.Minute)
						}
						break
					}
				}
				_ = helpers.SaveTasks(tasks)
			}
			return m, tea.Quit
		}
	}
	return m, nil
}
