//go:build linux

package capture

import (
	"encoding/binary"
	"errors"
	"sync"

	"github.com/jfreymuth/pulse"
)

type linuxCapture struct {
	client   *pulse.Client
	stream   *pulse.RecordStream
	mu       sync.Mutex
	closed   bool
	dataChan chan []byte
}

func newPlatformCapture() (AudioCapture, error) {
	client, err := pulse.NewClient()
	if err != nil {
		return nil, err
	}

	sinks, err := client.ListSinks()
	if err != nil {
		client.Close()
		return nil, err
	}

	if len(sinks) == 0 {
		client.Close()
		return nil, errors.New("no audio sink found")
	}

	c := &linuxCapture{
		client:   client,
		dataChan: make(chan []byte, 10),
	}

	writer := pulse.Int16Writer(func(buf []int16) (int, error) {
		c.mu.Lock()
		closed := c.closed
		c.mu.Unlock()

		if closed {
			return 0, errors.New("capture closed")
		}

		data := make([]byte, len(buf)*2)
		for i, sample := range buf {
			binary.LittleEndian.PutUint16(data[i*2:], uint16(sample))
		}

		select {
		case c.dataChan <- data:
		default:
		}

		return len(buf), nil
	})

	stream, err := client.NewRecord(
		writer,
		pulse.RecordMonitor(sinks[0]),
		pulse.RecordSampleRate(SampleRate),
		pulse.RecordStereo,
		pulse.RecordLatency(0.05),
	)
	if err != nil {
		client.Close()
		return nil, err
	}

	c.stream = stream

	return c, nil
}

func (c *linuxCapture) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return errors.New("capture is closed")
	}
	c.stream.Start()
	return nil
}

func (c *linuxCapture) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return nil
	}
	c.closed = true
	if c.stream != nil {
		c.stream.Close()
	}
	if c.client != nil {
		c.client.Close()
	}
	return nil
}

func (c *linuxCapture) Read(p []byte) (int, error) {
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

func (c *linuxCapture) SampleRate() int {
	return SampleRate
}

func (c *linuxCapture) Channels() int {
	return Channels
}

