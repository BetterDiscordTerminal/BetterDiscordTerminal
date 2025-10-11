# BDTerm

[![Release](https://img.shields.io/github/v/release/BetterDiscordTerminal/BetterDiscordTerminal)](https://github.com/BetterDiscordTerminal/BetterDiscordTerminal/releases/latest)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8.svg)](https://go.dev/)

A beautiful terminal-based installer for [BetterDiscord](https://betterdiscord.app) on macOS.

## Features

- **Interactive TUI** - Beautiful terminal interface built with [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **Install BetterDiscord** - Easily install BetterDiscord on Discord, PTB, or Canary
- **Uninstall** - Clean removal of BetterDiscord from any Discord version
- **Repair** - Fix broken BetterDiscord installations
- **Checksum Verification** - Ensures downloaded files are authentic
- **Apple Silicon Native** - Optimized for ARM64 Macs

## Installation

### Quick Install

Saves the executable in the /tmp directory.

```bash
curl -fsSL https://bd.pand.dev/install.sh | bash
```

```bash
curl -fsSL https://raw.githubusercontent.com/BetterDiscordTerminal/BetterDiscordTerminal/main/install.sh | bash
```

### Manual Install

1. Download the latest binary from [releases](https://github.com/BetterDiscordTerminal/BetterDiscordTerminal/releases/latest)
2. Make it executable:
   ```bash
   chmod +x bdterm
   ```
3. Move to your PATH:
   ```bash
   sudo mv bdterm /usr/local/bin/
   ```

## Usage

Simply run:

```bash
bdterm
```

Navigate the menu using:
- `↑/↓` - Move selection
- `Enter` - Select
- `ESC` - Go back
- `q` or `Ctrl+C` - Quit

## Requirements

- **macOS 11.0+** (Big Sur or later)
- **Apple Silicon (ARM64)** - Intel Macs are not currently supported
- **Discord** - Must have Discord, Discord PTB, or Discord Canary installed

## Supported Discord Versions

- **Discord** - Stable release
- **Discord PTB** - Public Test Build
- **Discord Canary** - Bleeding-edge experimental build

## Development

### Prerequisites

- Go 1.21 or later
- macOS with Apple Silicon

### Building from Source

```bash
# Clone the repository
git clone https://github.com/BetterDiscordTerminal/BetterDiscordTerminal.git
cd BetterDiscordTerminal

# Build
go build -o bdterm

# Run
./bdterm
```

### Build Script

```bash
# Build universal binary
./build.sh
```

## How It Works

1. **Detects Discord** - Finds installed Discord versions on your system
2. **Downloads BetterDiscord** - Fetches the latest release from GitHub
3. **Verifies Integrity** - Checks SHA256 checksums
4. **Injects Shim** - Safely modifies Discord to load BetterDiscord
5. **Launches Discord** - Automatically restarts Discord with BetterDiscord enabled

## Uninstallation

To remove BetterDiscord:

1. Run `bdterm`
2. Select "Uninstall"
3. Choose the Discord version to uninstall from
4. BetterDiscord will be cleanly removed

## Troubleshooting

### Installation Failed

- Make sure Discord is fully closed before installing
- Try running with `sudo` if you get permission errors
- Check that you have the latest version of Discord installed

### Discord Won't Launch

- Use the "Repair" option in BDTerm
- Manually verify Discord is installed in `/Applications`
- Reinstall Discord if issues persist

### BetterDiscord Not Loading

- Run the "Repair" option
- Check that BetterDiscord files exist in `~/Library/Application Support/BetterDiscord`
- Verify the shim was correctly injected

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Security

### Checksum Verification

BDTerm verifies all downloaded files using SHA256 checksums to ensure authenticity and prevent tampering.

### Source Code

All source code is open and available for inspection. The installer:
- Only modifies Discord's index.js to add BetterDiscord.asar module.
- Does not collect any data
- Does not make external network requests except to download BetterDiscord

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Disclaimer

This tool is not affiliated with Discord Inc. or BetterDiscord. Use at your own risk. Modifying Discord may violate their Terms of Service.

## Why this project was created

I wanted to make a TUI application using go and I wanted an installer that wasn't resource intensive.


Made with ❤️ by [nmn](https://pand.dev)
