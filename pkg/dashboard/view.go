package dashboard

import (
	"fmt"
	"io"
	"strconv"

	"github.com/austinemk/sigcat/pkg/helpers"
	"github.com/austinemk/sigcat/pkg/theme"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (d taskDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	t, ok := listItem.(helpers.BreakTask)
	if !ok {
		return
	}

	var taskRow string
	if index == m.Index() {
		// Highlighted active row layout using your theme definitions
		rawRow := fmt.Sprintf(
			" [%s] %s\n    Every %dm | %s",
			t.ID,
			t.Title,
			t.DurationMin,
			fmt.Sprintf("Ready (%s)", t.NextRun.Format("15:04:05")),
		)
		if !t.IsActive {
			rawRow = fmt.Sprintf(" [%s] %s\n    Every %dm | Inactive", t.ID, t.Title, t.DurationMin)
		}
		taskRow = theme.SelectedItemStyle.Render(rawRow)
	} else {
		// Default unselected row layout using your theme definitions
		activeStr := lipgloss.NewStyle().Foreground(theme.SubtleColor).Render("❌ Inactive")
		if t.IsActive {
			activeStr = lipgloss.NewStyle().Foreground(theme.GreenColor).Render(fmt.Sprintf("🟢 Ready (%s)", t.NextRun.Format("15:04:05")))
		}

		taskRow = fmt.Sprintf(
			"  %s %s\n    Every %s | %s",
			lipgloss.NewStyle().Foreground(theme.PurpleColor).Bold(true).Render("["+t.ID+"]"),
			t.Title,
			theme.FocusStyle.Render(strconv.Itoa(t.DurationMin)+"m"),
			activeStr,
		)
	}

	fmt.Fprint(w, taskRow)
}

func (m dashboardModel) View() tea.View {
	statusStr := lipgloss.NewStyle().Foreground(theme.RedColor).Render("[🔴 STOPPED]")
	if m.daemonRunning {
		statusStr = lipgloss.NewStyle().Foreground(theme.GreenColor).Render("[🟢 RUNNING]")
	}
	headerBar := theme.GenerateTexturedShadowTitle("SIGCAT HUB", "#AEB6FC", "#475569") + statusStr + "\n\n"

	var bodyContent string

	if m.state == viewTasks {
		bodyContent += "Active Automation Timers Matrix:\n\n"
		if len(m.taskList.Items()) == 0 {
			bodyContent += theme.HelpStyle.Render("  No active profiles found. Press [n] to create one.") + "\n"
		} else {
			bodyContent += m.taskList.View() + "\n"
		}

		bodyContent += "\n" + theme.HelpStyle.Render("[n] New Task • [space] Toggle • [s] Start/Stop Daemon • [d] Delete • [/] Filter • [q] Quit")
	} else {
		bodyContent += theme.TitleStyle.Render("✨ CREATE NEW SCHEDULER PROFILE") + "\n\n"

		labels := []string{"Window Title:   ", "Sweet Message:  ", "Timeout (Mins): ", "AutoRepeat(y/n):"}
		for i, label := range labels {
			rowText := fmt.Sprintf("  %s %s", label, m.inputs[i].View())
			if m.inputIndex == i {
				bodyContent += theme.SelectedItemStyle.Render(rowText) + "\n"
			} else {
				bodyContent += rowText + "\n"
			}
		}

		if m.errMessage != "" {
			bodyContent += "\n" + theme.ErrorStyle.Render("❌ "+m.errMessage) + "\n"
		}
		bodyContent += "\n" + theme.HelpStyle.Render("[Esc] Cancel • [Tab/Arrows] Navigate • [Enter] Next / Save")
	}

	boxedLayout := theme.CardStyle.Render(bodyContent)
	v := tea.NewView(headerBar + boxedLayout)
	v.AltScreen = true
	return v
}
