package tada

import (
	"errors"
	"fmt"
	"os"
)

func Main() {
	if err := Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func Run(args []string) error {
	path, err := TodoFilePath()
	if err != nil {
		return fmt.Errorf("resolving storage path: %w", err)
	}

	storage := NewStorage[Todos](path)
	todos := Todos{}

	if err := storage.Load(&todos); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("loading tasks: %w", err)
		}
	}

	if len(args) <= 1 || (len(args) == 2 && args[1] == "tui") {
		if err := RunTUI(&todos, storage); err != nil {
			return fmt.Errorf("TUI error: %w", err)
		}
		return nil
	}

	mutated, err := RunCommand(args[1:], &todos)
	if err != nil {
		return err
	}

	if mutated {
		if err := storage.Save(todos); err != nil {
			return fmt.Errorf("saving tasks: %w", err)
		}
	}
	return nil
}
