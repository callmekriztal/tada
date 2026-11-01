package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func printHelp() {
	fmt.Println("tada - A simple command-line task manager")
	fmt.Println("\nUsage:")
	fmt.Println("  tada                     Launch interactive TUI (default)")
	fmt.Println("  tada [command] [args]    Run CLI command")
	fmt.Println("  tada [flags]             Run with CLI flags")
	fmt.Println("\nCommands:")
	fmt.Println("  tui                      Launch interactive TUI")
	fmt.Println("  list, ls                 List all tasks in terminal table")
	fmt.Println("  add <title>              Add a new task")
	fmt.Println("  toggle, done <id>        Toggle task completion by ID")
	fmt.Println("  edit <id> <new title>    Edit a task title by ID")
	fmt.Println("  del, rm <id>             Delete a task by ID")
	fmt.Println("  clear [-f]               Clear all tasks (-f skips confirmation)")
	fmt.Println("  stats [-p <period>] [--json]   Show completion statistics (today/week/month/all)")
	fmt.Println("  help                     Show this help message")
	fmt.Println("\nFlags (backward compatibility):")
	fmt.Println("  --list                   List all tasks")
	fmt.Println("  --add <title>            Add a new task")
	fmt.Println("  --toggle <id>            Toggle task completion by ID")
	fmt.Println("  --edit <id:new_title>    Edit task by ID:new_title")
	fmt.Println("  --del <id>               Delete task by ID")
	fmt.Println("  --clear                  Clear all tasks")
	fmt.Println("  --stats                  Show completion statistics")
	fmt.Println("  --period, -p <period>    Filter statistics by period (today/week/month/all)")
	fmt.Println("  --json                   Output statistics in JSON format")
}

func parseIndex(str string) (int, error) {
	id, err := strconv.Atoi(str)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid task ID %q: must be a positive number", str)
	}
	return id - 1, nil
}

func confirmClear() bool {
	fmt.Print("Are you sure you want to clear all tasks? [y/N]: ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		ans := strings.TrimSpace(strings.ToLower(scanner.Text()))
		return ans == "y" || ans == "yes"
	}
	return false
}

func RunCommand(args []string, todos *Todos) (bool, error) {
	if len(args) == 0 {
		todos.print()
		return false, nil
	}

	first := args[0]

	// Handle flag-style commands (e.g. --add, -add, --list, etc.)
	if strings.HasPrefix(first, "-") {
		return runFlagCommand(args, todos)
	}

	// Handle subcommand-style commands
	switch strings.ToLower(first) {
	case "list", "ls":
		todos.print()
		return false, nil

	case "add":
		if len(args) < 2 {
			return false, errors.New("usage: tada add <task title>")
		}
		title := strings.TrimSpace(strings.Join(args[1:], " "))
		if title == "" {
			return false, errors.New("task title cannot be empty")
		}
		todos.add(title)
		fmt.Printf("Added task #%d: %q\n", len(*todos), title)
		return true, nil

	case "toggle", "done", "check":
		if len(args) < 2 {
			return false, errors.New("usage: tada done <task id>")
		}
		index, err := parseIndex(args[1])
		if err != nil {
			return false, err
		}
		if err := todos.toggle(index); err != nil {
			return false, err
		}
		item := (*todos)[index]
		if item.Completed {
			fmt.Printf("Completed task #%d: %q\n", index+1, item.Title)
		} else {
			fmt.Printf("Reopened task #%d: %q\n", index+1, item.Title)
		}
		return true, nil

	case "del", "delete", "rm", "remove":
		if len(args) < 2 {
			return false, errors.New("usage: tada del <task id>")
		}
		index, err := parseIndex(args[1])
		if err != nil {
			return false, err
		}
		if err := todos.validateIndex(index); err != nil {
			return false, err
		}
		title := (*todos)[index].Title
		if err := todos.delete(index); err != nil {
			return false, err
		}
		fmt.Printf("Deleted task #%d: %q\n", index+1, title)
		return true, nil

	case "edit":
		if len(args) < 2 {
			return false, errors.New("usage: tada edit <task id> <new title>")
		}
		var index int
		var newTitle string
		if strings.Contains(args[1], ":") {
			parts := strings.SplitN(args[1], ":", 2)
			var err error
			index, err = parseIndex(parts[0])
			if err != nil {
				return false, err
			}
			newTitle = parts[1]
			if len(args) > 2 {
				newTitle += " " + strings.Join(args[2:], " ")
			}
		} else {
			if len(args) < 3 {
				return false, errors.New("usage: tada edit <task id> <new title>")
			}
			var err error
			index, err = parseIndex(args[1])
			if err != nil {
				return false, err
			}
			newTitle = strings.Join(args[2:], " ")
		}

		newTitle = strings.TrimSpace(newTitle)
		if newTitle == "" {
			return false, errors.New("task title cannot be empty")
		}

		if err := todos.edit(index, newTitle); err != nil {
			return false, err
		}
		fmt.Printf("Updated task #%d: %q\n", index+1, newTitle)
		return true, nil

	case "clear":
		force := false
		if len(args) > 1 && (args[1] == "-f" || args[1] == "--force") {
			force = true
		}
		if !force && !confirmClear() {
			fmt.Println("Clear cancelled.")
			return false, nil
		}
		todos.clear()
		fmt.Println("All tasks have been cleared.")
		return true, nil

	case "stats", "stat", "statistics":
		return runStatsCommand(args[1:], todos)

	case "help":
		printHelp()
		return false, nil

	default:
		printHelp()
		return false, fmt.Errorf("unknown command: %q", first)
	}
}

func runStatsCommand(args []string, todos *Todos) (bool, error) {
	fs := flag.NewFlagSet("stats", flag.ContinueOnError)
	var periodFlag string
	var jsonOutput bool

	fs.StringVar(&periodFlag, "p", "", "Period filter (today/week/month/all)")
	fs.StringVar(&periodFlag, "period", "", "Period filter (today/week/month/all)")
	fs.BoolVar(&jsonOutput, "json", false, "Output in JSON format")

	var cleanFlags []string
	var posArgs []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "-p" || arg == "--period" {
			if i+1 < len(args) {
				cleanFlags = append(cleanFlags, arg, args[i+1])
				i++
			} else {
				cleanFlags = append(cleanFlags, arg)
			}
		} else if strings.HasPrefix(arg, "-") {
			cleanFlags = append(cleanFlags, arg)
		} else {
			posArgs = append(posArgs, arg)
		}
	}

	if err := fs.Parse(cleanFlags); err != nil {
		return false, err
	}

	period := periodFlag
	if period == "" && len(posArgs) > 0 {
		period = posArgs[0]
	}
	if period == "" {
		period = "all"
	}

	stats, err := ComputeStats(*todos, period, time.Now())
	if err != nil {
		return false, err
	}

	if err := OutputStats(stats, jsonOutput); err != nil {
		return false, err
	}
	return false, nil
}

func runFlagCommand(args []string, todos *Todos) (bool, error) {
	fs := flag.NewFlagSet("tada", flag.ContinueOnError)

	var (
		add         = fs.String("add", "", "Add a new Task")
		edit        = fs.String("edit", "", "Edit a task by id:new_title")
		del         = fs.Int("del", -1, "Delete a task by index")
		toggle      = fs.Int("toggle", -1, "Toggle task completion by index")
		list        = fs.Bool("list", false, "List all tasks")
		clear       = fs.Bool("clear", false, "Clear all tasks")
		force       = fs.Bool("force", false, "Force clear without prompt")
		stats       = fs.Bool("stats", false, "Show completion statistics")
		period      = fs.String("period", "", "Period filter for stats (today/week/month/all)")
		periodShort = fs.String("p", "", "Period filter for stats (today/week/month/all)")
		jsonOut     = fs.Bool("json", false, "Output in JSON format")
	)

	fs.Usage = printHelp

	if err := fs.Parse(args); err != nil {
		return false, err
	}

	switch {
	case *list:
		todos.print()
		return false, nil

	case *stats || *period != "" || *periodShort != "":
		p := *period
		if *periodShort != "" {
			p = *periodShort
		}
		if p == "" {
			p = "all"
		}
		s, err := ComputeStats(*todos, p, time.Now())
		if err != nil {
			return false, err
		}
		if err := OutputStats(s, *jsonOut); err != nil {
			return false, err
		}
		return false, nil

	case *add != "":
		title := strings.TrimSpace(*add)
		if title == "" {
			return false, errors.New("task title cannot be empty")
		}
		todos.add(title)
		fmt.Printf("Added task #%d: %q\n", len(*todos), title)
		return true, nil

	case *edit != "":
		parts := strings.SplitN(*edit, ":", 2)
		if len(parts) != 2 {
			return false, errors.New("invalid format for --edit, use id:new_title")
		}
		index, err := parseIndex(parts[0])
		if err != nil {
			return false, err
		}
		newTitle := strings.TrimSpace(parts[1])
		if newTitle == "" {
			return false, errors.New("task title cannot be empty")
		}
		if err := todos.edit(index, newTitle); err != nil {
			return false, err
		}
		fmt.Printf("Updated task #%d: %q\n", index+1, newTitle)
		return true, nil

	case *toggle != -1:
		index := *toggle - 1
		if err := todos.toggle(index); err != nil {
			return false, err
		}
		item := (*todos)[index]
		if item.Completed {
			fmt.Printf("Completed task #%d: %q\n", index+1, item.Title)
		} else {
			fmt.Printf("Reopened task #%d: %q\n", index+1, item.Title)
		}
		return true, nil

	case *del != -1:
		index := *del - 1
		if err := todos.validateIndex(index); err != nil {
			return false, err
		}
		title := (*todos)[index].Title
		if err := todos.delete(index); err != nil {
			return false, err
		}
		fmt.Printf("Deleted task #%d: %q\n", index+1, title)
		return true, nil

	case *clear:
		if !*force && !confirmClear() {
			fmt.Println("Clear cancelled.")
			return false, nil
		}
		todos.clear()
		fmt.Println("All tasks have been cleared.")
		return true, nil

	default:
		todos.print()
		return false, nil
	}
}
