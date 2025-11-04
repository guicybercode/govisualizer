package capture

import "io"

const (
	SampleRate = 44100
	Channels   = 2
	BitDepth   = 16
)

type AudioCapture interface {
	Start() error
	Stop() error
	Read([]byte) (int, error)
	SampleRate() int
	Channels() int
}

func NewCapture() (AudioCapture, error) {
	return newPlatformCapture()
}

type Reader struct {
	capture AudioCapture
}

func NewReader(capture AudioCapture) io.Reader {
	return &Reader{capture: capture}
}

func (r *Reader) Read(p []byte) (int, error) {
	return r.capture.Read(p)
}

