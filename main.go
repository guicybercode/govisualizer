package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/guicybercode/govisualizer/internal/capture"
	"github.com/guicybercode/govisualizer/internal/processor"
	"github.com/guicybercode/govisualizer/internal/ui"
)

func main() {
	audioCapture, err := capture.NewCapture()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize audio capture: %v\n", err)
		os.Exit(1)
	}

	err = audioCapture.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start audio capture: %v\n", err)
		os.Exit(1)
	}
	defer audioCapture.Stop()

	fftProc := processor.NewFFTProcessor(audioCapture.SampleRate())
	spectrogram := processor.NewSpectrogram(100, processor.FrequencyBands)

	go processAudio(audioCapture, fftProc, spectrogram)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		audioCapture.Stop()
		os.Exit(0)
	}()

	model := ui.NewModel(spectrogram)
	program := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}

func processAudio(capture capture.AudioCapture, fftProc *processor.FFTProcessor, spectrogram *processor.Spectrogram) {
	buffer := make([]byte, 4096)
	sampleBuffer := make([]int16, processor.FFTSize*capture.Channels())
	sampleIndex := 0

	for {
		n, err := capture.Read(buffer)
		if err != nil {
			break
		}

		if n == 0 {
			continue
		}

		for i := 0; i < n; i += 2 {
			if sampleIndex >= len(sampleBuffer) {
				mono := make([]int16, processor.FFTSize)
				for j := 0; j < processor.FFTSize; j++ {
					if capture.Channels() == 2 {
						left := sampleBuffer[j*2]
						right := sampleBuffer[j*2+1]
						mono[j] = (left + right) / 2
					} else {
						mono[j] = sampleBuffer[j]
					}
				}

				magnitudes := fftProc.Process(mono)
				if magnitudes != nil {
					normalized := fftProc.Normalize(magnitudes)
					spectrogram.AddLine(normalized)
				}

				copy(sampleBuffer[:sampleIndex], sampleBuffer[sampleIndex:])
				sampleIndex = 0
			}

			if i+1 < n {
				val := int16(binary.LittleEndian.Uint16(buffer[i : i+2]))
				sampleBuffer[sampleIndex] = val
				sampleIndex++
			}
		}
	}
}

