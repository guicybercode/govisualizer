package processor

import (
	"sync"
)

type Spectrogram struct {
	data     [][]float64
	mu       sync.RWMutex
	maxLines int
	bands    int
}

func NewSpectrogram(maxLines, bands int) *Spectrogram {
	return &Spectrogram{
		data:     make([][]float64, 0, maxLines),
		maxLines: maxLines,
		bands:    bands,
	}
}

func (s *Spectrogram) AddLine(magnitudes []float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(magnitudes) != s.bands {
		return
	}

	line := make([]float64, s.bands)
	copy(line, magnitudes)

	s.data = append(s.data, line)

	if len(s.data) > s.maxLines {
		s.data = s.data[1:]
	}
}

func (s *Spectrogram) GetData() [][]float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([][]float64, len(s.data))
	for i, line := range s.data {
		result[i] = make([]float64, len(line))
		copy(result[i], line)
	}

	return result
}

func (s *Spectrogram) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = s.data[:0]
}

func (s *Spectrogram) Lines() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}

