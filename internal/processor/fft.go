package processor

import (
	"math"

	"github.com/mjibson/go-dsp/fft"
)

const (
	FFTSize      = 2048
	FrequencyBands = 128
	MaxFrequency = 20000
)

type FFTProcessor struct {
	window     []float64
	sampleRate int
}

func NewFFTProcessor(sampleRate int) *FFTProcessor {
	window := make([]float64, FFTSize)
	for i := 0; i < FFTSize; i++ {
		window[i] = 0.5 * (1 - math.Cos(2*math.Pi*float64(i)/float64(FFTSize-1)))
	}
	return &FFTProcessor{
		window:     window,
		sampleRate: sampleRate,
	}
}

func (p *FFTProcessor) Process(samples []int16) []float64 {
	if len(samples) < FFTSize {
		return nil
	}

	complexSamples := make([]complex128, FFTSize)
	for i := 0; i < FFTSize; i++ {
		val := float64(samples[i]) / 32768.0
		val *= p.window[i]
		complexSamples[i] = complex(val, 0)
	}

	fftResult := fft.FFT(complexSamples)

	magnitudes := make([]float64, FrequencyBands)
	binWidth := float64(p.sampleRate) / float64(FFTSize)
	
	for i := 0; i < FrequencyBands; i++ {
		targetFreq := float64(i) * float64(MaxFrequency) / float64(FrequencyBands)
		binIndex := int(targetFreq / binWidth)
		
		if binIndex >= len(fftResult) {
			binIndex = len(fftResult) - 1
		}
		
		mag := math.Sqrt(real(fftResult[binIndex])*real(fftResult[binIndex]) + 
			imag(fftResult[binIndex])*imag(fftResult[binIndex]))
		
		magnitudes[i] = mag
	}

	return magnitudes
}

func (p *FFTProcessor) Normalize(magnitudes []float64) []float64 {
	if len(magnitudes) == 0 {
		return magnitudes
	}

	max := 0.0
	for _, m := range magnitudes {
		if m > max {
			max = m
		}
	}

	if max == 0 {
		return magnitudes
	}

	normalized := make([]float64, len(magnitudes))
	for i, m := range magnitudes {
		normalized[i] = m / max
	}

	return normalized
}

func (p *FFTProcessor) ToDB(magnitudes []float64) []float64 {
	db := make([]float64, len(magnitudes))
	for i, m := range magnitudes {
		if m <= 0 {
			db[i] = -100
		} else {
			db[i] = 20 * math.Log10(m)
		}
	}
	return db
}

