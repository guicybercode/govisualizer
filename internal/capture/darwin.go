//go:build darwin

package capture

import (
	"errors"
	"sync"

	"github.com/gordonklaus/portaudio"
)

type darwinCapture struct {
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

	devices, err := portaudio.Devices()
	if err != nil {
		portaudio.Terminate()
		return nil, err
	}

	var device *portaudio.Device
	for _, d := range devices {
		if d.MaxInputChannels > 0 && (d.Name == "BlackHole 2ch" || d.Name == "BlackHole 16ch" || d.Name == "Loopback Audio") {
			device = d
			break
		}
	}

	if device == nil {
		device = host.DefaultInputDevice
	}

	if device == nil {
		portaudio.Terminate()
		return nil, errors.New("no input device found")
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

	c := &darwinCapture{
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

func (c *darwinCapture) callback(in []int16) {
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

func (c *darwinCapture) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return errors.New("capture is closed")
	}
	return c.stream.Start()
}

func (c *darwinCapture) Stop() error {
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

func (c *darwinCapture) Read(p []byte) (int, error) {
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

func (c *darwinCapture) SampleRate() int {
	return SampleRate
}

func (c *darwinCapture) Channels() int {
	return Channels
}

