# GoTodo

A simple command-line Todo List application written in Go.

GoTodo lets you manage tasks directly from your terminal and persists them locally using JSON.

## Features

* **Zero external dependencies**: Built 100% on the Go standard library.
* **Fast & lightweight**: Clean, instant CLI feedback on all actions.
* **Intuitive subcommands & flags**: Support for both `todocli add` and `todocli --add`.
* **Safe operations**: Confirmation prompt when clearing tasks (bypassable with `-f`).
* **Atomic JSON persistence**: Prevents file corruption on unexpected interruptions.
* Add, list, edit, toggle, and delete tasks.
* Compact timestamps and formatted terminal table.
* Full unit test suite with 100% standard library tests.

## Storage

Tasks are stored locally at:

```text
~/.todocli/todos.json
```

No database or external service is required.

## Installation

### Linux — AppImage

Download `GoTodo-x86_64.AppImage` from the GitHub Releases page.

Make it executable:

```bash
chmod +x GoTodo-x86_64.AppImage
```

Run it:

```bash
./GoTodo-x86_64.AppImage
```

The AppImage does not require installation or Go to be installed.

### Linux — Standalone Binary

Download `todocli-linux-amd64` from the GitHub Releases page.

Make it executable:

```bash
chmod +x todocli-linux-amd64
```

For system-wide installation, move it to `/usr/local/bin`:

```bash
sudo install -m 755 todocli-linux-amd64 /usr/local/bin/todocli
```

You can now run it from anywhere:

```bash
todocli
```

Check that it is available:

```bash
which todocli
```

Expected:

```text
/usr/local/bin/todocli
```

### Build from Source

Requirements:

* Go 1.20 or newer
* Linux, macOS, or Windows

Clone the repository and build:

```bash
git clone https://github.com/callmekriztal/Go_Todo.git
cd Go_Todo
go build -trimpath -ldflags="-s -w" -o todocli .
```

Install system-wide on Linux:

```bash
sudo install -m 755 todocli /usr/local/bin/todocli
```

Verify:

```bash
todocli help
```

## Usage

### List tasks (default)

Simply run `todocli` or `todocli list`:

```bash
todocli
# or
todocli list
```

### Add a task

```bash
todocli add "Learn Go"
```

### Toggle task completion

```bash
todocli done 1
# or
todocli toggle 1
```

Running the command again reopens the task.

### Edit a task

```bash
todocli edit 1 "Learn advanced Go"
```

### Delete a task

```bash
todocli del 1
# or
todocli rm 1
```

### Clear all tasks

Prompts for confirmation to prevent accidental loss:

```bash
todocli clear
```

To skip the confirmation prompt:

```bash
todocli clear -f
```

### Flags compatibility

All commands are also accessible via standard flags:

```bash
todocli --list
todocli --add "Learn Go"
todocli --toggle 1
todocli --edit 1:"Learn advanced Go"
todocli --del 1
todocli --clear
```

### Show help

```bash
todocli help
# or
todocli --help
```

## Example

```bash
$ todocli add "Learn Go"
✓ Added task #1: "Learn Go"

$ todocli add "Build a CLI app"
✓ Added task #2: "Build a CLI app"

$ todocli done 1
✓ Completed task #1: "Learn Go"

$ todocli
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
├── command.go       # CLI flags and command execution
├── helper.go        # Todo operations and display
├── storage.go       # JSON persistence
├── main.go          # Application entry point
├── helper_test.go   # Unit tests
├── go.mod
├── go.sum
└── README.md
```

## Release

Prebuilt Linux binaries are available through GitHub Releases:

* `GoTodo-x86_64.AppImage` — portable Linux AppImage
* `todocli-linux-amd64` — standalone Linux binary

The AppImage can be run without installation. The standalone binary can be installed system-wide using `/usr/local/bin`.

## License

MIT License

Made by Christy Dominic Cyrl
