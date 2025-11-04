//go:build windows

package capture

import (
	"errors"
	"sync"

	"github.com/gordonklaus/portaudio"
)

type windowsCapture struct {
	stream   *portaudio.Stream
	mu       sync.Mutex
	dataChan chan []byte
	closed   bool
}

func newPlatformCapture() (AudioCapture, error) {
	if err := portaudio.Initialize(); err != nil {
		return nil, err
	}

	host, err := portaudio.DefaultHostApi()
	if err != nil {
		portaudio.Terminate()
		return nil, err
	}

	device := host.DefaultOutputDevice
	if device == nil {
		portaudio.Terminate()
		return nil, errors.New("no output device found")
	}

	framesPerBuffer := 1024

	params := portaudio.StreamParameters{
		Input: portaudio.StreamDeviceParameters{
			Device:  device,
			Latency: device.DefaultLowInputLatency,
		},
		Output:         nil,
		SampleRate:     SampleRate,
		FramesPerBuffer: framesPerBuffer,
		Flags:          portaudio.ClipOff,
	}

	c := &windowsCapture{
		dataChan: make(chan []byte, 10),
	}

	stream, err := portaudio.OpenStream(params, c.callback)
	if err != nil {
		portaudio.Terminate()
		return nil, err
	}

	c.stream = stream

	return c, nil
}

func (c *windowsCapture) callback(in []int16) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return
	}

	data := make([]byte, len(in)*2)
	for i, sample := range in {
		data[i*2] = byte(sample)
		data[i*2+1] = byte(sample >> 8)
	}

	select {
	case c.dataChan <- data:
	default:
	}
}

func (c *windowsCapture) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return errors.New("capture is closed")
	}
	return c.stream.Start()
}

func (c *windowsCapture) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return nil
	}
	c.closed = true
	err := c.stream.Stop()
	if err != nil {
		return err
	}
	err = c.stream.Close()
	if err != nil {
		return err
	}
	close(c.dataChan)
	return portaudio.Terminate()
}

func (c *windowsCapture) Read(p []byte) (int, error) {
	data, ok := <-c.dataChan
	if !ok {
		return 0, errors.New("capture stream closed")
	}

	n := copy(p, data)
	if n < len(data) {
		remaining := data[n:]
		select {
		case c.dataChan <- remaining:
		default:
		}
	}

	return n, nil
}

func (c *windowsCapture) SampleRate() int {
	return SampleRate
}

func (c *windowsCapture) Channels() int {
	return Channels
}

