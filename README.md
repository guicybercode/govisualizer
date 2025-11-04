# 🎵 Go Audio Spectrogram Visualizer

A real-time audio spectrogram visualizer that captures system audio playback and displays a beautiful waterfall-style frequency spectrum in the terminal using [Bubbletea](https://github.com/charmbracelet/bubbletea).

![Spectrogram Visualization](https://img.shields.io/badge/spectrogram-waterfall-blue)
![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-green)

## Features

- 🎨 **Beautiful Waterfall Visualization**: Real-time spectrogram with vibrant color gradients
- 🎧 **System Audio Capture**: Captures audio playback (not microphone) across platforms
- ⌨️ **Interactive Controls**: Pause/resume and quit functionality
- 🚀 **Cross-Platform**: Supports Windows, Linux, and macOS
- 📊 **Real-time FFT Processing**: Fast Fourier Transform for frequency analysis
- 🎯 **Modular Architecture**: Clean, separated packages for capture, processing, and UI

## Platform Support

### Windows
Uses WASAPI loopback mode via PortAudio to capture system audio output.

### Linux
Uses PulseAudio monitor source to capture audio from the default sink.

### macOS
Supports BlackHole or Loopback Audio virtual devices. You may need to install BlackHole:
```bash
brew install blackhole-2ch
```

## Installation

### Prerequisites

- Go 1.21 or higher
- For Linux: PulseAudio development libraries
- For macOS: BlackHole (optional, for loopback capture)

### Build from Source

```bash
git clone https://github.com/guicybercode/govisualizer.git
cd govisualizer
go mod download
go build -o govisualizer
```

### Run

```bash
./govisualizer
```

## Usage

1. Start the application - it will automatically begin capturing system audio
2. Play audio on your system (music, videos, etc.)
3. Watch the real-time spectrogram visualization

### Controls

- `Space` - Pause/Resume visualization
- `Q` or `Ctrl+C` - Quit the application

## Project Structure

```
govisualizer/
├── main.go                 # Entry point and orchestration
├── internal/
│   ├── capture/           # Platform-specific audio capture
│   │   ├── capture.go     # Interface and common logic
│   │   ├── windows.go     # WASAPI loopback (via PortAudio)
│   │   ├── linux.go       # PulseAudio monitor
│   │   └── darwin.go       # BlackHole/PortAudio
│   ├── processor/         # Audio processing and FFT
│   │   ├── fft.go         # FFT wrapper and frequency analysis
│   │   └── spectrogram.go # Spectrogram data structure
│   └── ui/                # Bubbletea TUI components
│       ├── model.go       # Main Bubbletea model
│       ├── view.go        # Spectrogram rendering
│       └── styles.go      # Color schemes and styling
└── README.md
```

## Technical Details

- **Sample Rate**: 44.1 kHz
- **FFT Size**: 2048 samples
- **Frequency Bands**: 128 bands (0-20 kHz)
- **Update Rate**: ~60 FPS
- **Audio Format**: 16-bit PCM, stereo

## Dependencies

- [Bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [PortAudio](https://github.com/gordonklaus/portaudio) - Audio I/O for Windows/macOS
- [Pulse](https://github.com/jfreymuth/pulse) - PulseAudio bindings for Linux
- [go-dsp](https://github.com/mjibson/go-dsp) - FFT implementation

## Educational Purpose

This project is designed for learning about:
- Digital audio processing
- Fast Fourier Transform (FFT)
- Real-time data visualization
- Cross-platform audio capture
- Terminal User Interface (TUI) development

## License

MIT License - feel free to use this project for learning and experimentation.

## Contributing

Contributions are welcome! This is an educational project, so feel free to:
- Report issues
- Suggest improvements
- Submit pull requests
- Share your modifications

---

**하나님은 나의 빛이요, 나의 구원이시니 내가 누구를 두려워하랴. (시편 27:1)**

*God is my light and my salvation—whom shall I fear? (Psalm 27:1)*

