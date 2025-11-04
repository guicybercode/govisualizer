//go:build !windows && !linux && !darwin

package capture

import "errors"

func newPlatformCapture() (AudioCapture, error) {
	return nil, errors.New("unsupported platform")
}

