package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/guicybercode/govisualizer/internal/processor"
)

type state int

const (
	stateRunning state = iota
	statePaused
	stateQuit
)

type Model struct {
	spectrogram *processor.Spectrogram
	state       state
	width       int
	height      int
}

func NewModel(spectrogram *processor.Spectrogram) Model {
	return Model{
		spectrogram: spectrogram,
		state:       stateRunning,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, tickCmd())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.state = stateQuit
			return m, tea.Quit
		case " ":
			if m.state == stateRunning {
				m.state = statePaused
			} else {
				m.state = stateRunning
			}
			return m, nil
		}

	case tickMsg:
		if m.state == stateRunning {
			return m, tickCmd()
		}
		return m, nil
	}

	return m, nil
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second/60, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

