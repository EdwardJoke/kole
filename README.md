# kole

A terminal-based environment variable manager built with Go and the Charm libraries.

![Go Version](https://img.shields.io/badge/Go-1.25.5-blue)
![License](https://img.shields.io/github/license/EdwardJoke/kole)

## Features

- **View All Environment Variables** - Browse through all environment variables in your current session
- **Manage PATH Variable** - View, delete, and reorder PATH entries
- **Add New Environment Variables** - Create new environment variables with validation
- **Edit Environment Variables** - Modify existing environment variable values
- **Delete Environment Variables** - Remove environment variables from your session
- **Search Environment Variables** - Filter environment variables by name or value
- **Format Shell Config** - Organize and format .bashrc and .zshrc files

## Installation

### Download Pre-built Binaries

Download the latest release from the [GitHub Releases](https://github.com/EdwardJoke/kole/releases) page.

### Build from Source

```bash
git clone https://github.com/EdwardJoke/kole.git
cd kole
go build
```

## Usage

Run the application:

```bash
./kole
```

Use the arrow keys to navigate through the menu options and press Enter to select.

## Building

This project uses [GoReleaser](https://goreleaser.com/) for building and releasing binaries.

```bash
goreleaser release --snapshot
```

## Dependencies

- [huh](https://github.com/charmbracelet/huh) - Interactive TUI forms
- [lipgloss](https://github.com/charmbracelet/lipgloss) - Style definition library

## License

Apache-2.0 License - see [LICENSE](LICENSE) for details.

## Thanks

- [charm](https://charm.sh/) - For the libraries used in this project