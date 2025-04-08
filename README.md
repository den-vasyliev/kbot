# kbot

A Telegram bot for controlling traffic light signals using GPIO pins on a Raspberry Pi.

## Features

- Control traffic light signals (red, amber, green) through Telegram commands
- Toggle individual lights on/off
- Simple and intuitive command interface
- GPIO pin control for Raspberry Pi
- Cross-platform support (Linux, Darwin, Windows)
- Multi-architecture support (amd64, arm64)
- Docker containerization support

## Prerequisites

- Raspberry Pi with GPIO access
- Go 1.16 or later
- Docker (optional, for containerization)
- Telegram Bot Token (set as TELE_TOKEN environment variable)
- Required Go packages:
  - github.com/spf13/cobra
  - github.com/stianeikeland/go-rpio
  - gopkg.in/telebot.v4

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/kbot.git
cd kbot
```

2. Set up your Telegram Bot Token:
```bash
export TELE_TOKEN="your_telegram_bot_token"
```

3. Build the application:
```bash
make build
```

## Build Options

The project supports various build targets through Makefile:

- `make format` - Format Go code
- `make lint` - Run golint
- `make test` - Run tests
- `make get` - Get dependencies
- `make build` - Build the application
- `make image` - Build Docker image
- `make push` - Push Docker image to registry
- `make clean` - Clean build artifacts

### Build Configuration

You can customize the build by setting environment variables:
```bash
TARGETOS=linux    # Target OS (linux, darwin, windows)
TARGETARCH=arm64  # Target architecture (amd64, arm64)
```

## Docker Deployment

1. Build the Docker image:
```bash
make image
```

2. Push to registry (optional):
```bash
make push
```

The Docker image will be tagged as:
```
${REGISTRY}/${APP}:${VERSION}-${TARGETARCH}
```

## Usage

Start the bot:
```bash
./kbot start
```

### Available Commands

- `/s red` - Toggle red light
- `/s amber` - Toggle amber light
- `/s green` - Toggle green light
- `hello` - Get a greeting from the bot

### GPIO Pin Configuration

The bot uses the following GPIO pins by default:
- Red light: GPIO 12
- Amber light: GPIO 27
- Green light: GPIO 22

## Development

The project uses Cobra for CLI command management and go-rpio for GPIO control.

### Project Structure

- `cmd/` - Contains the main command implementations
  - `kbot.go` - Main bot implementation and traffic light control
  - `root.go` - Root command configuration
  - `version.go` - Version command implementation

## Versioning

The application version is automatically generated during build using:
- Latest git tag
- Short commit hash

## License

This project is licensed under the MIT License - see the LICENSE file for details.
