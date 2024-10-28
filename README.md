# k-assist

k-assist (command: `kass`) is a cross-platform terminal assistant focused on simplicity. It helps developers by providing AI-powered assistance for terminal commands and code-related queries.

## Features

- AI-powered command suggestions, chat, and error handling
- Context, directory, and shell history aware
- Cross-platform support (Linux, macOS, Windows)
- Support for multiple shells (bash, zsh, powershell)

## Notes

- You can generate your free gemini api key [here](https://ai.google.dev/gemini-api/docs/api-key)

- Don't forget to edit the configuration file and add your api key after installation. More details [here](#configuration).

- Kass does not preserve a chat session with the LLM, you won't need to worry about your previous messages affecting the current one. However, kass does have access to your shell history so it will see your previous commands when needed.

- It is not recommended to use the `-a` and `-A` flag inside big directories like home/, as it may cause unexpected errors due to the possibility of it containing sensitive data, and violating LLM providers' usage policies.

- If you encounter a safety error from the LLM provider that doesn't make sense, it is most likely due to the above.

- Usage of `-a` and `-A` will substantially increase the time it takes to generate a response depending on the size of the directory. You may need to adjust your `max_tokens` setting in the configuration file if you want to use them extensively.

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

For linux and macOS:

```bash
make build
make install
```

For windows:

```bash
The windows installer is not yet available. :(
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

kass can output multiple commands when necessary, line by line. Each command can be edited and executed individually. Use enter to execute the command, or ^C to interrupt the output.

### Chat Functionality

To enable chat functionality, use the `-c` flag:

```bash
kass -c "explain how can i create a recovery image for my system"
```

### Including All Directory Contents

To include all subdirectories in the context, use the `-a` flag:

```bash
kass -a "compress all the pdf files in the project docs directory"
```

### Including All Directory Contents with Data

To include all subdirectories and files with their contents, use the `-A` flag:

```bash
kass -A "create a compressed backup of the project excluding node_modules, build directories, and temporary files""
```

### Combining Flags

You can combine multiple flags to get content aware assistance easily.

```bash
kass -A -c "analyze my project structure and suggest improvements"
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
