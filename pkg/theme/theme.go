package theme

import (
	"charm.land/lipgloss/v2"
)

// ==========================================
// 1. Color Palette Tokens
// ==========================================
var (
	SubtleColor = lipgloss.Color("#64748B")
	PurpleColor = lipgloss.Color("#AEB6FC")
	PinkColor   = lipgloss.Color("#FFB8D1")
	GreenColor  = lipgloss.Color("#22C55E")
	RedColor    = lipgloss.Color("#EF4444")
	DarkSlate   = lipgloss.Color("#1E293B")
)

// ==========================================
// 2. Structural & Text Styles
// ==========================================
var (
	TitleStyle = lipgloss.NewStyle().Foreground(PurpleColor).Bold(true).Padding(0, 1)
	CardStyle  = lipgloss.NewStyle().Padding(1, 2).MarginBottom(1)
	ErrorStyle = lipgloss.NewStyle().Foreground(RedColor).Bold(true)
	HelpStyle  = lipgloss.NewStyle().Foreground(SubtleColor).Italic(true)
	FocusStyle = lipgloss.NewStyle().Foreground(PinkColor).Bold(true)

	// Clean background highlight style for active selections
	SelectedItemStyle = lipgloss.NewStyle().
				Background(PurpleColor).
				Foreground(DarkSlate).
				Bold(true)
)
