package main

import (
	_ "embed" // Required for go:embed directive compiler hooks
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/common-nighthawk/go-figure"
)

// 1. Instruct the Go compiler to bake bloody.flf into this byte slice
//
//go:embed Bloody.flf
var bloodyFontBytes []byte

// GenerateTexturedShadowTitle handles fonts like "Bloody" embedded directly in the binary.
func GenerateTexturedShadowTitle(input string, fgColor, shadowColor string) string {
	fgStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(fgColor))
	shadowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(shadowColor))

	// 2. Convert the embedded byte slice into an io.Reader on the fly
	fontReader := strings.NewReader(string(bloodyFontBytes))

	// 3. Pass the reader directly into go-figure. No errors, no disk files required!
	myFigure := figure.NewFigureWithFont(input, fontReader, true)
	rawAscii := myFigure.String()

	// 4. Run your exact drop-shadow tracing map logic
	lines := strings.Split(rawAscii, "\n")
	var completedBanner []string

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		runes := []rune(line)
		var builtLine strings.Builder
		shouldShadow := false

		for i := 0; i < len(runes); i++ {
			char := runes[i]
			isTextureCharacter := char != ' ' && char != '\u00a0'

			if isTextureCharacter {
				builtLine.WriteString(fgStyle.Render(string(char)))
				shouldShadow = true
			} else if char == ' ' && shouldShadow {
				builtLine.WriteString(shadowStyle.Render("│"))
				shouldShadow = false
			} else {
				builtLine.WriteRune(' ')
				shouldShadow = false
			}
		}
		completedBanner = append(completedBanner, builtLine.String())
	}

	return strings.Join(completedBanner, "\n")
}

// TruncateString cuts a string to a specified max length.
// If addEllipsis is true and the string is truncated, it appends "..."
func TruncateString(s string, maxLength int, addEllipsis bool) string {
	// Clean up leading/trailing whitespace
	s = strings.TrimSpace(s)

	// Convert to runes to handle UTF-8/Unicode characters safely
	runes := []rune(s)

	// If the string is already short enough, return it as-is
	if len(runes) <= maxLength {
		return s
	}

	// Truncate to the maximum allowed characters
	truncated := string(runes[:maxLength])

	// Trim any trailing space created by the cut before adding dots
	if addEllipsis {
		truncated = strings.TrimSpace(truncated) + "..."
	}

	return truncated
}
