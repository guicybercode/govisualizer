package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	var sb strings.Builder

	header := headerStyle.Render("🎵 Audio Spectrogram Visualizer 🎵")
	sb.WriteString(header)
	sb.WriteString("\n\n")

	spectrogramHeight := m.height - 6
	if spectrogramHeight < 1 {
		spectrogramHeight = 1
	}

	data := m.spectrogram.GetData()
	spectrogramWidth := m.width - 2

	if len(data) == 0 {
		empty := strings.Repeat(" ", spectrogramWidth)
		for i := 0; i < spectrogramHeight; i++ {
			sb.WriteString(empty)
			sb.WriteString("\n")
		}
	} else {
		startLine := 0
		if len(data) > spectrogramHeight {
			startLine = len(data) - spectrogramHeight
		}

		for i := startLine; i < len(data); i++ {
			line := data[i]
			sb.WriteString(renderSpectrogramLine(line, spectrogramWidth))
			sb.WriteString("\n")
		}

		for i := len(data) - startLine; i < spectrogramHeight; i++ {
			sb.WriteString(strings.Repeat(" ", spectrogramWidth))
			sb.WriteString("\n")
		}
	}

	sb.WriteString("\n")

	status := "● Running"
	statusColor := lipgloss.Color("#00ff00")
	if m.state == statePaused {
		status = "⏸ Paused"
		statusColor = lipgloss.Color("#ffcc00")
	}

	statusText := lipgloss.NewStyle().Foreground(statusColor).Render(status)
	controls := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Render(
		"Space: Pause/Resume | Q: Quit",
	)

	footer := lipgloss.JoinHorizontal(lipgloss.Left, statusText, "  ", controls)
	sb.WriteString(footerStyle.Render(footer))

	return sb.String()
}

func renderSpectrogramLine(magnitudes []float64, width int) string {
	if len(magnitudes) == 0 {
		return strings.Repeat(" ", width)
	}

	var sb strings.Builder
	bandsPerChar := float64(len(magnitudes)) / float64(width)

	for i := 0; i < width; i++ {
		startBand := int(float64(i) * bandsPerChar)
		endBand := int(float64(i+1) * bandsPerChar)
		if endBand > len(magnitudes) {
			endBand = len(magnitudes)
		}
		if startBand >= endBand {
			startBand = endBand - 1
		}
		if startBand < 0 {
			startBand = 0
		}

		avgIntensity := 0.0
		count := 0
		for j := startBand; j < endBand && j < len(magnitudes); j++ {
			avgIntensity += magnitudes[j]
			count++
		}
		if count > 0 {
			avgIntensity /= float64(count)
		}

		char := intensityToChar(avgIntensity)
		color := frequencyToColor(startBand, len(magnitudes))
		styled := styleText(char, color)
		sb.WriteString(styled)
	}

	return sb.String()
}

