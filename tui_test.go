package main

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestTUIInitAndRender(t *testing.T) {
	todos := Todos{}
	todos.add("Task 1")
	todos.add("Task 2")

	m := newTUIModel(&todos, nil)
	m.width = 80
	m.height = 24

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
	todos := Todos{}
	todos.add("Task 1")
	todos.add("Task 2")

	m := newTUIModel(&todos, nil)

	// Move down
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updated.(tuiModel)
	if m.cursor != 1 {
		t.Fatalf("expected cursor at 1, got %d", m.cursor)
	}

	// Move up
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = updated.(tuiModel)
	if m.cursor != 0 {
		t.Fatalf("expected cursor at 0, got %d", m.cursor)
	}

	// Toggle Task 1
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeySpace})
	m = updated.(tuiModel)
	if !todos[0].Completed {
		t.Fatal("expected Task 1 to be completed after Space")
	}

	// Toggle Task 1 back
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(tuiModel)
	if todos[0].Completed {
		t.Fatal("expected Task 1 to be incomplete after Enter")
	}
}

func TestTUIFilterTabs(t *testing.T) {
	todos := Todos{}
	todos.add("Task 1")
	todos.add("Task 2")
	_ = todos.toggle(1) // Task 2 completed

	m := newTUIModel(&todos, nil)
	if len(m.filteredIndices) != 2 {
		t.Fatalf("expected 2 tasks in FilterAll, got %d", len(m.filteredIndices))
	}

	// Switch to FilterActive
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(tuiModel)
	if m.filter != FilterActive {
		t.Fatalf("expected FilterActive, got %v", m.filter)
	}
	if len(m.filteredIndices) != 1 || m.filteredIndices[0] != 0 {
		t.Fatalf("expected 1 active task, got %d", len(m.filteredIndices))
	}

	// Switch to FilterCompleted
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(tuiModel)
	if m.filter != FilterCompleted {
		t.Fatalf("expected FilterCompleted, got %v", m.filter)
	}
	if len(m.filteredIndices) != 1 || m.filteredIndices[0] != 1 {
		t.Fatalf("expected 1 completed task, got %d", len(m.filteredIndices))
	}

	// Switch back to FilterAll
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(tuiModel)
	if m.filter != FilterAll {
		t.Fatalf("expected FilterAll, got %v", m.filter)
	}
	if len(m.filteredIndices) != 2 {
		t.Fatalf("expected 2 tasks in FilterAll, got %d", len(m.filteredIndices))
	}
}

func TestTUIAddTaskInline(t *testing.T) {
	todos := Todos{}
	m := newTUIModel(&todos, nil)

	// Press 'a' to begin adding
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	m = updated.(tuiModel)
	if !m.adding {
		t.Fatal("expected adding state to be true")
	}

	// Type characters
	for _, r := range "New TUI task" {
		updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		m = updated.(tuiModel)
	}

	// Press Enter to submit
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(tuiModel)
	if m.adding {
		t.Fatal("expected adding state to be false after Enter")
	}
	if len(todos) != 1 || todos[0].Title != "New TUI task" {
		t.Fatalf("expected 1 task 'New TUI task', got %v", todos)
	}

	// Test Esc to cancel
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	m = updated.(tuiModel)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = updated.(tuiModel)
	if m.adding {
		t.Fatal("expected adding state to be false after Esc")
	}
	if len(todos) != 1 {
		t.Fatalf("expected 1 task still, got %d", len(todos))
	}
}

func TestTUIDeleteTask(t *testing.T) {
	todos := Todos{}
	todos.add("Task 1")
	todos.add("Task 2")

	m := newTUIModel(&todos, nil)

	// Delete Task 1
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	m = updated.(tuiModel)

	if len(todos) != 1 || todos[0].Title != "Task 2" {
		t.Fatalf("expected only Task 2 remaining, got %v", todos)
	}
	if len(m.filteredIndices) != 1 {
		t.Fatalf("expected 1 filtered index, got %d", len(m.filteredIndices))
	}
}

func TestTUIQuit(t *testing.T) {
	todos := Todos{}
	m := newTUIModel(&todos, nil)

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil {
		t.Fatal("expected quit command on 'q'")
	}
}
