# tada

A simple, fast command-line task manager written in Go.

**tada** makes managing tasks quick and straightforward right from your terminal, with a sleek interactive TUI, comprehensive statistics, and local JSON persistence.

## Features

* **Interactive TUI**: Modern, centered terminal user interface built with Charm's `bubbletea` and `lipgloss`.
* **Statistics & Streaks**: View completion metrics, daily averages, streaks, and `#tag` breakdowns across today, week, month, or all-time periods (`tada stats`).
* **Quick Task Entry**: Press `a` or `n` in the TUI to add tasks inline.
* **Instant Filtering**: Press `Tab` to cycle between All, Active, and Completed views.
* **Seamless CLI & TUI integration**: Run `tada` for the TUI, or use `tada add`, `tada stats`, etc. from scripts.
* **Safe operations**: Confirmation prompt when clearing tasks (bypassable with `-f`).
* **Atomic JSON persistence**: Prevents file corruption on unexpected interruptions.
* Full unit test suite with high coverage.

## Storage

Tasks are stored locally at:

```text
~/.tada/todos.json
```

No database or external service is required. Seamlessly migrates existing tasks from `~/.todocli/todos.json`.

## Usage

### Interactive TUI (default)

Simply run `tada` with no arguments to launch the interactive terminal interface:

```bash
tada
```

#### TUI Keyboard Shortcuts:
| Key | Action |
| --- | --- |
| `↑` / `k`, `↓` / `j` | Navigate tasks up / down |
| `Space` / `Enter` | Toggle task completion |
| `a` / `n` | Add task (inline input, `Enter` to save, `Esc` to cancel) |
| `d` / `x` | Delete selected task |
| `Tab` / `Shift+Tab` | Cycle filter: **All** → **Active** → **Completed** |
| `q` / `Ctrl+C` | Quit TUI cleanly |

### CLI Commands

You can also use `tada` directly from the command line:

```bash
# List tasks in table format
tada list

# Add a task
tada add "Learn Go"

# Complete/toggle a task
tada done 1

# Edit a task
tada edit 1 "Learn advanced Go"

# Delete a task
tada del 1

# Clear all tasks
tada clear
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

### Statistics & Streaks

View completion metrics, completion rate, daily average, current/longest streaks, and `#tag` breakdowns:

```bash
# View all-time statistics
tada stats

# Filter by time period: today, week (7 days), month (30 days), all
tada stats today
tada stats week
tada stats month
tada stats all

# Use the -p or --period flag
tada stats -p week
tada stats --period month

# Output in JSON format (e.g. for scripting or integrations)
tada stats --json
tada stats -p week --json
```

Example statistics output:

```text
Todo Statistics (Last 7 days)
────────────────────────────────────────
Tasks Completed:     8
Tasks Created:       10
Completion Rate:     80.0%
Daily Average:       2.5 tasks/day
Current Streak:      3 day(s)
Longest Streak:      7 day(s)
────────────────────────────────────────
Tag Breakdown:
  #cli:              5 completed
  #work:             3 completed
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
tada --stats
tada --stats -p week
tada --stats --json
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
tada/
├── command.go       # CLI commands, subcommands, and flags
├── helper.go        # Todo operations, display, and path helpers
├── storage.go       # Atomic JSON persistence
├── stats.go         # Statistics computation, periods, streaks, and formatting
├── tui.go           # Interactive terminal UI (Bubble Tea & Lipgloss)
├── main.go          # Application entry point
├── helper_test.go   # Helper and command tests
├── stats_test.go    # Statistics and period range tests
├── tui_test.go      # TUI state and navigation tests
├── go.mod
├── go.sum
└── README.md
```

## Release

Prebuilt Linux binaries are available through [GitHub Releases](https://github.com/callmekriztal/tada/releases):

* `tada-linux-amd64` — standalone Linux binary
* `tada-x86_64.AppImage` — portable Linux AppImage

## License

MIT License

Made by Christy Dominic Cyrl

