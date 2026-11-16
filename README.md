# tada

A fast, delightful command-line task manager written in Go.

**tada** delivers effortless task management directly in your terminal, featuring a modern interactive Terminal UI (TUI), comprehensive productivity statistics with streak tracking, and atomic local JSON persistence.

---

## Features

- **Interactive TUI**: Centered, distraction-free interface powered by Bubble Tea & Lipgloss.
- **Productivity Statistics**: Track completion rates, daily averages, streaks, and `#tag` breakdowns across periods (`today`, `week`, `month`, `all`).
- **Quick Task Capture**: Add tasks inline in the TUI (`a` / `n`) or instantly via CLI.
- **Smart Filtering**: Switch between All, Active, and Completed views in seconds (`Tab`).
- **Unified CLI & TUI**: Launch the TUI by running `tada`, or use scripts and shell commands (`tada add`, `tada stats --json`).
- **Atomic Persistence**: Prevents file corruption via temporary file swaps in `~/.tada/todos.json`.
- **System-Wide Ready**: Installable via Go, single-line binary download, or portable AppImage.

---

## Installation

### 1. One-Line System-Wide Install (Linux / macOS)

Download the latest prebuilt binary directly into `/usr/local/bin`:

```bash
curl -fsSL https://github.com/callmekriztal/tada/releases/latest/download/tada-linux-amd64 -o tada
chmod +x tada
sudo install -m 755 tada /usr/local/bin/tada
rm tada
```

### 2. Via `go install`

If you have Go installed (`go 1.20+`):

```bash
go install github.com/callmekriztal/tada/v2@latest
```

Ensure `$(go env GOPATH)/bin` is in your `$PATH`.

### 3. Build from Source

```bash
git clone https://github.com/callmekriztal/tada.git
cd tada
go build -trimpath -ldflags="-s -w" -o tada .
sudo install -m 755 tada /usr/local/bin/tada
```

### 4. Portable AppImage

```bash
curl -fsSL https://github.com/callmekriztal/tada/releases/latest/download/tada-x86_64.AppImage -o tada.AppImage
chmod +x tada.AppImage
./tada.AppImage
```

All releases and binary artifacts are available on [GitHub Releases](https://github.com/callmekriztal/tada/releases).

---

## Usage

### Interactive TUI (Default)

Running `tada` without arguments launches the interactive interface:

```bash
tada
```

#### Keyboard Shortcuts

| Key | Action |
| --- | --- |
| `↑` / `k`, `↓` / `j` | Navigate tasks up / down |
| `Space` / `Enter` | Toggle task completion |
| `a` / `n` | Add task (inline input, `Enter` to save, `Esc` to cancel) |
| `d` / `x` | Delete selected task |
| `Tab` / `Shift+Tab` | Cycle filter: **All** → **Active** → **Completed** |
| `q` / `Ctrl+C` | Quit TUI |

---

### CLI Commands

Use `tada` directly from scripts or your shell:

```bash
# List tasks
tada list                # or 'tada ls'

# Add tasks (quotes optional for multi-word tasks)
tada add Buy groceries for dinner
tada add Review PR #work #urgent

# Toggle task completion (run again to reopen)
tada done 1              # or 'tada toggle 1'

# Edit task title
tada edit 1 Update documentation

# Delete task by ID
tada del 1               # or 'tada rm 1'

# Clear all tasks (prompts for confirmation)
tada clear
tada clear -f            # skip prompt
```

#### Flags Compatibility

Standard flags are supported for full backward compatibility:

```bash
tada --list
tada --add "Write tests"
tada --toggle 1
tada --edit 1:"Updated title"
tada --del 1
tada --clear
tada --stats
tada --stats -p week
tada --stats --json
```

---

## Statistics & Streaks

Track your productivity and habits across custom time windows:

```bash
# All-time statistics
tada stats

# Filter by period: today, week (7d), month (30d), all
tada stats today
tada stats week
tada stats month

# Period flag syntax
tada stats -p week
tada stats --period month

# Export structured JSON (ideal for status bars, scripts, or backups)
tada stats --json
tada stats -p week --json
```

#### Sample Output

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
  #urgent:           2 completed
  #work:             6 completed
```

- **Streaks**: A streak continues if yesterday was completed, giving you until the end of today to log a completion.
- **Tags**: Any `#tag` in task titles is automatically tracked and aggregated.

---

## Storage & Migration

Tasks are stored locally as JSON at:

```text
~/.tada/todos.json
```

- **Zero configuration**: No database daemon or network service required.
- **Safe writes**: Writes to a temporary file (`.tmp`) and swaps atomically.
- **Auto-migration**: Automatically detects and migrates legacy tasks from `~/.todocli/todos.json`.

---

## Development

Run tests:

```bash
go test ./...
```

Run tests with verbose output:

```bash
go test -v ./...
```

Run static analysis:

```bash
go vet ./...
```

---

## Project Structure

```text
tada/
├── main.go          # Application entry point
├── pkg/
│   └── tada/        # Core application package
│       ├── app.go       # App lifecycle & router
│       ├── command.go   # Subcommands & flag parsing
│       ├── helper.go    # Task models, table display & path resolvers
│       ├── stats.go     # Statistics, streaks & tag breakdown
│       ├── storage.go   # Atomic JSON persistence
│       └── tui.go       # Terminal UI (Bubble Tea & Lipgloss)
├── tests/           # Dedicated test suite
│       ├── helper_test.go
│       ├── stats_test.go
│       └── tui_test.go
├── go.mod
├── go.sum
└── README.md
```

---

## License

MIT License © Christy Dominic Cyril
