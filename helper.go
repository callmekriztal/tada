package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"
	"time"
)

var ErrInvalidIndex = errors.New("invalid index")

type Todo struct {
	Title       string     `json:"title"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

type Todos []Todo

// Done is kept as a type alias for backward compatibility.
type Done = Todos

func (t *Todos) add(title string) {
	todo := Todo{
		Title:       title,
		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
	}
	*t = append(*t, todo)
}

func (t *Todos) validateIndex(index int) error {
	if index < 0 || index >= len(*t) {
		return ErrInvalidIndex
	}
	return nil
}

func (t *Todos) delete(index int) error {
	if err := t.validateIndex(index); err != nil {
		return err
	}

	*t = append((*t)[:index], (*t)[index+1:]...)
	return nil
}

func (t *Todos) toggle(index int) error {
	if err := t.validateIndex(index); err != nil {
		return err
	}

	item := &(*t)[index]
	item.Completed = !item.Completed
	if item.Completed {
		completionTime := time.Now()
		item.CompletedAt = &completionTime
	} else {
		item.CompletedAt = nil
	}

	return nil
}

func (t *Todos) edit(index int, title string) error {
	if err := t.validateIndex(index); err != nil {
		return err
	}

	(*t)[index].Title = title
	return nil
}

func (t *Todos) print() {
	if len(*t) == 0 {
		fmt.Println("No tasks found! Run 'tada add \"task description\"' to add one.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	fmt.Fprintln(w, "ID\tSTATUS\tTITLE\tCREATED\tCOMPLETED")

	for index, item := range *t {
		status := "[ ]"
		completedAt := "-"

		if item.Completed {
			status = "[✓]"
			if item.CompletedAt != nil {
				completedAt = item.CompletedAt.Format("2006-01-02 15:04")
			}
		}

		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
			index+1,
			status,
			item.Title,
			item.CreatedAt.Format("2006-01-02 15:04"),
			completedAt,
		)
	}
	w.Flush()
}

func todoFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".tada")

	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}

	target := filepath.Join(dir, "todos.json")

	// Seamless migration from old ~/.todocli/todos.json if needed
	if _, err := os.Stat(target); errors.Is(err, os.ErrNotExist) {
		oldFile := filepath.Join(home, ".todocli", "todos.json")
		if data, err := os.ReadFile(oldFile); err == nil {
			_ = os.WriteFile(target, data, 0644)
		}
	}

	return target, nil
}

func (t *Todos) clear() {
	*t = Todos{}
}
