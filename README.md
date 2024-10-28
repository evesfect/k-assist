# k-assist

k-assist (command: `kass`) is a cross-platform terminal assistant focused on simplicity. It helps developers by providing AI-powered assistance for terminal commands and code-related queries.

## Features

- AI-powered command suggestions, chat, and error handling
- Context, directory, and shell history aware
- Cross-platform support (Linux, macOS, Windows)
- Support for multiple shells (bash, zsh, powershell)

## Installation

### Prerequisites

- Go 1.23.2 or higher
- Make (for building from source)

### Building from Source

1. Clone the repository:

```bash
git clone https://github.com/evesfect/k-assist.git
cd k-assist
```

1. Build and install:

```bash
make build
make cleaninstall
```

## Configuration

k-assist uses a configuration file to determine custom user settings.
The configuration file is automatically created at:

- Linux/macOS: `~/.config/kass/config.json`
- Windows: `%APPDATA%\kass\config.json`

Example configuration:

```json
{
    "os": "arch-linux",
    "user": "evesfect",
    "llm": {
        "provider": "gemini",
        "api_key": "your-gemini-api-key-here",
        "model": "gemini-pro"
    },
    "max_tokens": 50,
    "shell": "bash"
}
```

## Usage

### Basic Command Assistance

kass is used to provide command assistance by default.

```bash
kass "how many docker containers are running right now?"

docker ps | wc -l
```

### Chat Functionality

To enable chat functionality, use the `-c` flag:

```bash
kass -c "Explain how can i create a recovery image for my system"
```

### Including All Directory Contents

To include all subdirectories and files in the context, use the `-a` flag:

```bash
kass -a "compress all the pdf files in the project docs directory"
```

### Combining Flags

You can combine multiple flags to get content aware assistance easily.

```bash
kass -a -c "Analyze my project structure and suggest improvements"
```

### Error Assistance

When you encounter an error, k-assist can help troubleshoot it:

1. If something fails along the way, k-assist will offer to help
2. Type 'Y' to get assistance with the error

Kass will have access to the error message, command output, directory contents, and shell history to provide better assistance.

## Uninstallation

To uninstall k-assist:

```bash
make uninstall
```

## License

This project is licensed under the GNU General Public License v3.0 - see the LICENSE file for details.
