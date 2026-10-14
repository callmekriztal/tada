package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	path, err := todoFilePath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving storage path: %v\n", err)
		os.Exit(1)
	}

	storage := NewStorage[Todos](path)
	todos := Todos{}

	if err := storage.Load(&todos); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(os.Stderr, "Error loading tasks: %v\n", err)
			os.Exit(1)
		}
	}

	mutated, err := RunCommand(os.Args[1:], &todos)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if mutated {
		if err := storage.Save(todos); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving tasks: %v\n", err)
			os.Exit(1)
		}
	}
}
