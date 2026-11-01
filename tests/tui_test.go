package tests

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"tada/pkg/tada"
)

func TestTUIInitAndRender(t *testing.T) {
	todos := tada.Todos{}
	todos.Add("Task 1")
	todos.Add("Task 2")

	m := tada.NewTUIModel(&todos, nil)
	m.SetSize(80, 24)

	view := m.View()
	if !strings.Contains(view, "TADA TODO") {
		t.Fatalf("expected view to contain header, got:\n%s", view)
	}
	if !strings.Contains(view, "Task 1") || !strings.Contains(view, "Task 2") {
		t.Fatalf("expected view to contain tasks, got:\n%s", view)
	}
	if !strings.Contains(view, "All (2)") {
		t.Fatalf("expected view to contain tab stats, got:\n%s", view)
	}
}

func TestTUINavigationAndToggle(t *testing.T) {
	todos := tada.Todos{}
	todos.Add("Task 1")
	todos.Add("Task 2")

	m := tada.NewTUIModel(&todos, nil)

	// Move down
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updated.(tada.TUIModel)
	if m.Cursor() != 1 {
		t.Fatalf("expected cursor at 1, got %d", m.Cursor())
	}

	// Move up
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = updated.(tada.TUIModel)
	if m.Cursor() != 0 {
		t.Fatalf("expected cursor at 0, got %d", m.Cursor())
	}

	// Toggle Task 1
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeySpace})
	m = updated.(tada.TUIModel)
	if !todos[0].Completed {
		t.Fatal("expected Task 1 to be completed after Space")
	}

	// Toggle Task 1 back
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(tada.TUIModel)
	if todos[0].Completed {
		t.Fatal("expected Task 1 to be incomplete after Enter")
	}
}

func TestTUIFilterTabs(t *testing.T) {
	todos := tada.Todos{}
	todos.Add("Task 1")
	todos.Add("Task 2")
	_ = todos.Toggle(1) // Task 2 completed

	m := tada.NewTUIModel(&todos, nil)
	if len(m.FilteredIndices()) != 2 {
		t.Fatalf("expected 2 tasks in FilterAll, got %d", len(m.FilteredIndices()))
	}

	// Switch to FilterActive
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(tada.TUIModel)
	if m.Filter() != tada.FilterActive {
		t.Fatalf("expected FilterActive, got %v", m.Filter())
	}
	if len(m.FilteredIndices()) != 1 || m.FilteredIndices()[0] != 0 {
		t.Fatalf("expected 1 active task, got %d", len(m.FilteredIndices()))
	}

	// Switch to FilterCompleted
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(tada.TUIModel)
	if m.Filter() != tada.FilterCompleted {
		t.Fatalf("expected FilterCompleted, got %v", m.Filter())
	}
	if len(m.FilteredIndices()) != 1 || m.FilteredIndices()[0] != 1 {
		t.Fatalf("expected 1 completed task, got %d", len(m.FilteredIndices()))
	}

	// Switch back to FilterAll
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(tada.TUIModel)
	if m.Filter() != tada.FilterAll {
		t.Fatalf("expected FilterAll, got %v", m.Filter())
	}
	if len(m.FilteredIndices()) != 2 {
		t.Fatalf("expected 2 tasks in FilterAll, got %d", len(m.FilteredIndices()))
	}
}

func TestTUIAddTaskInline(t *testing.T) {
	todos := tada.Todos{}
	m := tada.NewTUIModel(&todos, nil)

	// Press 'a' to begin adding
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	m = updated.(tada.TUIModel)
	if !m.Adding() {
		t.Fatal("expected adding state to be true")
	}

	// Type characters
	for _, r := range "New TUI task" {
		updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		m = updated.(tada.TUIModel)
	}

	// Press Enter to submit
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(tada.TUIModel)
	if m.Adding() {
		t.Fatal("expected adding state to be false after Enter")
	}
	if len(todos) != 1 || todos[0].Title != "New TUI task" {
		t.Fatalf("expected 1 task 'New TUI task', got %v", todos)
	}

	// Test Esc to cancel
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	m = updated.(tada.TUIModel)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = updated.(tada.TUIModel)
	if m.Adding() {
		t.Fatal("expected adding state to be false after Esc")
	}
	if len(todos) != 1 {
		t.Fatalf("expected 1 task still, got %d", len(todos))
	}
}

func TestTUIDeleteTask(t *testing.T) {
	todos := tada.Todos{}
	todos.Add("Task 1")
	todos.Add("Task 2")

	m := tada.NewTUIModel(&todos, nil)

	// Delete Task 1
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	m = updated.(tada.TUIModel)

	if len(todos) != 1 || todos[0].Title != "Task 2" {
		t.Fatalf("expected only Task 2 remaining, got %v", todos)
	}
	if len(m.FilteredIndices()) != 1 {
		t.Fatalf("expected 1 filtered index, got %d", len(m.FilteredIndices()))
	}
}

func TestTUIQuit(t *testing.T) {
	todos := tada.Todos{}
	m := tada.NewTUIModel(&todos, nil)

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil {
		t.Fatal("expected quit command on 'q'")
	}
}
