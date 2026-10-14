# tada

A simple, fast command-line task manager written in Go.

**tada** makes managing tasks quick and straightforward right from your terminal, with zero external dependencies and local JSON persistence.

## Features

* **Zero external dependencies**: Built 100% on the Go standard library.
* **Instant feedback**: Clear confirmation on every task operation.
* **Intuitive subcommands & flags**: Support for both `tada add` and `tada --add`.
* **Safe operations**: Confirmation prompt when clearing tasks (bypassable with `-f`).
* **Atomic JSON persistence**: Prevents file corruption on unexpected interruptions.
* Add, list, edit, toggle, and delete tasks.
* Compact timestamps and formatted terminal table.
* Full unit test suite with 100% standard library tests.

## Storage

Tasks are stored locally at:

```text
~/.tada/todos.json
```

No database or external service is required. Seamlessly migrates existing tasks from `~/.todocli/todos.json`.

## Installation

### Linux — Standalone Binary

Download `tada-linux-amd64` from the GitHub Releases page.

Make it executable and install:

```bash
chmod +x tada-linux-amd64
sudo install -m 755 tada-linux-amd64 /usr/local/bin/tada
```

You can now run it from anywhere:

```bash
tada
```

Check that it is available:

```bash
which tada
```

Expected:

```text
/usr/local/bin/tada
```

### Linux — AppImage

Download `tada-x86_64.AppImage` from the GitHub Releases page.

Make it executable:

```bash
chmod +x tada-x86_64.AppImage
./tada-x86_64.AppImage
```

The AppImage does not require installation or Go to be installed.

### Build from Source

Requirements:

* Go 1.20 or newer
* Linux, macOS, or Windows

Clone the repository and build:

```bash
git clone https://github.com/callmekriztal/Go_Todo.git
cd Go_Todo
go build -trimpath -ldflags="-s -w" -o tada .
```

Install system-wide on Linux:

```bash
sudo install -m 755 tada /usr/local/bin/tada
```

Verify:

```bash
tada help
```

## Usage

### List tasks (default)

Simply run `tada` or `tada list`:

```bash
tada
# or
tada list
```

### Add a task

```bash
tada add "Learn Go"
```

### Complete a task

```bash
tada done 1
# or
tada toggle 1
```

Running the command again reopens the task.

### Edit a task

```bash
tada edit 1 "Learn advanced Go"
```

### Delete a task

```bash
tada del 1
# or
tada rm 1
```

### Clear all tasks

Prompts for confirmation to prevent accidental loss:

```bash
tada clear
```

To skip the confirmation prompt:

```bash
tada clear -f
```

### Flags compatibility

All commands are also accessible via standard flags:

```bash
tada --list
tada --add "Learn Go"
tada --toggle 1
tada --edit 1:"Learn advanced Go"
tada --del 1
tada --clear
```

### Show help

```bash
tada help
# or
tada --help
```

## Example

```bash
$ tada add "Learn Go"
Added task #1: "Learn Go"

$ tada add "Build a CLI app"
Added task #2: "Build a CLI app"

$ tada done 1
Completed task #1: "Learn Go"

$ tada
ID    STATUS    TITLE              CREATED             COMPLETED
1     [✓]       Learn Go           2026-10-14 23:30    2026-10-14 23:31
2     [ ]       Build a CLI app    2026-10-14 23:30    -
```

## Testing

Run the complete test suite:

```bash
go test ./...
```

Run tests with verbose output:

```bash
go test -v
```

Run static analysis:

```bash
go vet ./...
```

## Project Structure

```text
Go_Todo/
├── command.go       # CLI commands, subcommands, and flags
├── helper.go        # Todo operations, display, and path helpers
├── storage.go       # Atomic JSON persistence
├── main.go          # Application entry point
├── helper_test.go   # Unit tests
├── go.mod
├── go.sum
└── README.md
```

## Release

Prebuilt Linux binaries are available through GitHub Releases:

* `tada-linux-amd64` — standalone Linux binary
* `tada-x86_64.AppImage` — portable Linux AppImage

## License

MIT License

Made by Christy Dominic Cyrl

