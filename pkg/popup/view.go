package popup

import (
	"github.com/austinemk/sigcat/pkg/helpers"
	"github.com/austinemk/sigcat/pkg/theme"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// In view.go, update your View() method:

func (m PopupModel) View() tea.View {
	// 1. Safely retrieve the active Braille string frame
	activeArt := ""
	if len(m.frames) > 0 {
		activeArt = m.frames[m.currentFrame]
	}

	// 2. Render your layout elements (same as you had before)
	artStyled := lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Bold(true).Render(activeArt)
	bannerTitle := theme.GenerateTexturedShadowTitle(theme.TruncateString(m.Task.Title, 10, false), "#AEB6FC")

	var panel string
	panel += theme.AccentStyle.Width(30).Render(m.Task.Message) + "\n\n"
	panel += theme.MutedStyle.Render("status: "+helpers.Ternary(m.DaemonRunning, "● daemon active", "● daemon stopped")) + "\n"

	if m.Task.AutoRepeat {
		panel += theme.ActiveStye.Render("mode: AutoRepeat\n")
		panel += theme.MutedStyle.Render("[r] stop repeat [s/q] quit")
	} else {
		panel += theme.ActiveStye.Render("mode: Run Once\n")
		panel += theme.MutedStyle.Render("[r] repeat [s/q] quit")
	}

	// 3. Join layout elements together
	uiLayout := lipgloss.JoinHorizontal(lipgloss.Bottom, panel, artStyled)
	combinedView := lipgloss.JoinVertical(lipgloss.Center, bannerTitle, uiLayout)

	// 4. DYNAMICALLY CENTER THE CONTENT
	// If we haven't received the dimensions yet, fallback to a clean string
	var finalRender string
	if m.TerminalWidth > 0 && m.TerminalHeight > 0 {
		finalRender = lipgloss.Place(
			m.TerminalWidth,
			m.TerminalHeight,
			lipgloss.Center,
			lipgloss.Center,
			combinedView,
		)
	} else {
		// Fallback padding just in case the initial size event is slightly delayed
		finalRender = lipgloss.NewStyle().Padding(2, 4).Render(combinedView)
	}

	v := tea.NewView(finalRender)
	v.AltScreen = true
	return v
}
