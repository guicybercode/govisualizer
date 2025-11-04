package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	gradientColors = []string{
		"#000000",
		"#1a0d33",
		"#330066",
		"#4d0099",
		"#6600cc",
		"#8000ff",
		"#0066ff",
		"#00ccff",
		"#00ffcc",
		"#00ff66",
		"#33ff00",
		"#ccff00",
		"#ffcc00",
		"#ff9900",
		"#ff6600",
		"#ff3300",
		"#ff0000",
		"#ffffff",
	}
)

func frequencyToColor(freqIndex, maxFreq int) lipgloss.Color {
	if maxFreq == 0 {
		return lipgloss.Color("#000000")
	}

	ratio := float64(freqIndex) / float64(maxFreq)
	colorIndex := int(ratio * float64(len(gradientColors)-1))
	
	if colorIndex >= len(gradientColors) {
		colorIndex = len(gradientColors) - 1
	}
	
	return lipgloss.Color(gradientColors[colorIndex])
}

func intensityToChar(intensity float64) string {
	if intensity < 0.1 {
		return " "
	} else if intensity < 0.2 {
		return "░"
	} else if intensity < 0.4 {
		return "▒"
	} else if intensity < 0.6 {
		return "▓"
	} else {
		return "█"
	}
}

func styleText(text string, color lipgloss.Color) string {
	return lipgloss.NewStyle().Foreground(color).Render(text)
}

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00ffcc")).
			Align(lipgloss.Center).
			Padding(0, 1)

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Align(lipgloss.Center).
			Padding(0, 1)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00ff00")).
			Bold(true)
)

