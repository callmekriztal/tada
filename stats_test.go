package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func TestPeriodRangeAndPeriodRangeAt(t *testing.T) {
	refTime := time.Date(2026, 11, 1, 14, 30, 0, 0, time.Local)

	// Test periodRange wrapper
	from, to := periodRange("today")
	if from.IsZero() || to.IsZero() {
		t.Fatalf("periodRange('today') returned zero times")
	}

	// 1. Today
	from, to, label, err := periodRangeAt("today", refTime)
	if err != nil {
		t.Fatalf("unexpected error for today: %v", err)
	}
	expectedFrom := time.Date(2026, 11, 1, 0, 0, 0, 0, time.Local)
	if !from.Equal(expectedFrom) {
		t.Fatalf("expected from %v, got %v", expectedFrom, from)
	}
	if !to.Equal(refTime) {
		t.Fatalf("expected to %v, got %v", refTime, to)
	}
	if label != "Today" {
		t.Fatalf("expected label 'Today', got %q", label)
	}

	// 2. Week (7 days)
	from, to, label, err = periodRangeAt("week", refTime)
	if err != nil {
		t.Fatalf("unexpected error for week: %v", err)
	}
	expectedWeekFrom := refTime.AddDate(0, 0, -7)
	if !from.Equal(expectedWeekFrom) {
		t.Fatalf("expected from %v, got %v", expectedWeekFrom, from)
	}
	if label != "Last 7 days" {
		t.Fatalf("expected label 'Last 7 days', got %q", label)
	}

	// 3. Month (30 days)
	from, to, label, err = periodRangeAt("month", refTime)
	if err != nil {
		t.Fatalf("unexpected error for month: %v", err)
	}
	expectedMonthFrom := refTime.AddDate(0, 0, -30)
	if !from.Equal(expectedMonthFrom) {
		t.Fatalf("expected from %v, got %v", expectedMonthFrom, from)
	}
	if label != "Last 30 days" {
		t.Fatalf("expected label 'Last 30 days', got %q", label)
	}

	// 4. All
	from, to, label, err = periodRangeAt("all", refTime)
	if err != nil {
		t.Fatalf("unexpected error for all: %v", err)
	}
	if !from.IsZero() {
		t.Fatalf("expected zero from for all, got %v", from)
	}
	if label != "All time" {
		t.Fatalf("expected label 'All time', got %q", label)
	}

	// 5. Invalid
	_, _, _, err = periodRangeAt("invalid_period", refTime)
	if err == nil {
		t.Fatalf("expected error for invalid period, got nil")
	}
}

func TestComputeStats_EmptyTodos(t *testing.T) {
	now := time.Date(2026, 11, 1, 12, 0, 0, 0, time.Local)
	todos := Todos{}

	stats, err := ComputeStats(todos, "all", now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stats.TasksCompleted != 0 || stats.TasksCreated != 0 {
		t.Fatalf("expected 0 tasks, got completed=%d, created=%d", stats.TasksCompleted, stats.TasksCreated)
	}
	if stats.CompletionRate != 0 {
		t.Fatalf("expected 0 completion rate, got %f", stats.CompletionRate)
	}
	if stats.DailyAverage != 0 {
		t.Fatalf("expected 0 daily average, got %f", stats.DailyAverage)
	}
	if stats.CurrentStreak != 0 || stats.LongestStreak != 0 {
		t.Fatalf("expected 0 streaks, got current=%d, longest=%d", stats.CurrentStreak, stats.LongestStreak)
	}
	if len(stats.TagBreakdown) != 0 {
		t.Fatalf("expected empty tag breakdown, got %v", stats.TagBreakdown)
	}
}

func TestComputeStats_MetricsAndFiltering(t *testing.T) {
	now := time.Date(2026, 11, 1, 12, 0, 0, 0, time.Local)
	tToday := time.Date(2026, 11, 1, 9, 0, 0, 0, time.Local)
	t3DaysAgo := time.Date(2026, 10, 29, 10, 0, 0, 0, time.Local)
	t10DaysAgo := time.Date(2026, 10, 22, 11, 0, 0, 0, time.Local)
	t40DaysAgo := time.Date(2026, 9, 22, 8, 0, 0, 0, time.Local)

	todos := Todos{
		// Created & completed today
		{Title: "Task today #work", Completed: true, CreatedAt: tToday, CompletedAt: &tToday},
		// Created 3 days ago, completed 3 days ago
		{Title: "Task 3d #work #urgent", Completed: true, CreatedAt: t3DaysAgo, CompletedAt: &t3DaysAgo},
		// Created 10 days ago, completed 10 days ago
		{Title: "Task 10d #home", Completed: true, CreatedAt: t10DaysAgo, CompletedAt: &t10DaysAgo},
		// Created 40 days ago, completed 40 days ago
		{Title: "Task 40d", Completed: true, CreatedAt: t40DaysAgo, CompletedAt: &t40DaysAgo},
		// Created today, uncompleted
		{Title: "Unfinished today #work", Completed: false, CreatedAt: tToday, CompletedAt: nil},
		// Created 10 days ago, uncompleted
		{Title: "Unfinished 10d", Completed: false, CreatedAt: t10DaysAgo, CompletedAt: nil},
	}

	// 1. Period: today
	statsToday, err := ComputeStats(todos, "today", now)
	if err != nil {
		t.Fatalf("ComputeStats today failed: %v", err)
	}
	if statsToday.TasksCompleted != 1 {
		t.Fatalf("today: expected 1 completed, got %d", statsToday.TasksCompleted)
	}
	if statsToday.TasksCreated != 2 {
		t.Fatalf("today: expected 2 created, got %d", statsToday.TasksCreated)
	}
	if statsToday.CompletionRate != 50.0 {
		t.Fatalf("today: expected 50.0%% completion rate, got %.1f", statsToday.CompletionRate)
	}
	if statsToday.DailyAverage != 1.0 {
		t.Fatalf("today: expected 1.0 daily average, got %.1f", statsToday.DailyAverage)
	}
	if statsToday.TagBreakdown["#work"] != 1 {
		t.Fatalf("today: expected 1 #work, got %d", statsToday.TagBreakdown["#work"])
	}

	// 2. Period: week (last 7 days includes today & 3 days ago)
	statsWeek, err := ComputeStats(todos, "week", now)
	if err != nil {
		t.Fatalf("ComputeStats week failed: %v", err)
	}
	if statsWeek.TasksCompleted != 2 {
		t.Fatalf("week: expected 2 completed, got %d", statsWeek.TasksCompleted)
	}
	if statsWeek.TasksCreated != 3 {
		t.Fatalf("week: expected 3 created, got %d", statsWeek.TasksCreated)
	}
	// 2 / 3 * 100 = 66.7%
	if statsWeek.CompletionRate != 66.7 {
		t.Fatalf("week: expected 66.7%% completion rate, got %.1f", statsWeek.CompletionRate)
	}
	// 2 completed across 2 active days = 1.0 tasks/day
	if statsWeek.DailyAverage != 1.0 {
		t.Fatalf("week: expected 1.0 daily average, got %.1f", statsWeek.DailyAverage)
	}
	if statsWeek.TagBreakdown["#work"] != 2 || statsWeek.TagBreakdown["#urgent"] != 1 {
		t.Fatalf("week: tag breakdown mismatch: %v", statsWeek.TagBreakdown)
	}

	// 3. Period: month (last 30 days includes today, 3d, 10d)
	statsMonth, err := ComputeStats(todos, "month", now)
	if err != nil {
		t.Fatalf("ComputeStats month failed: %v", err)
	}
	if statsMonth.TasksCompleted != 3 {
		t.Fatalf("month: expected 3 completed, got %d", statsMonth.TasksCompleted)
	}
	if statsMonth.TasksCreated != 5 {
		t.Fatalf("month: expected 5 created, got %d", statsMonth.TasksCreated)
	}
	// 3 / 5 * 100 = 60.0%
	if statsMonth.CompletionRate != 60.0 {
		t.Fatalf("month: expected 60.0%% completion rate, got %.1f", statsMonth.CompletionRate)
	}

	// 4. Period: all
	statsAll, err := ComputeStats(todos, "all", now)
	if err != nil {
		t.Fatalf("ComputeStats all failed: %v", err)
	}
	if statsAll.TasksCompleted != 4 {
		t.Fatalf("all: expected 4 completed, got %d", statsAll.TasksCompleted)
	}
	if statsAll.TasksCreated != 6 {
		t.Fatalf("all: expected 6 created, got %d", statsAll.TasksCreated)
	}
	// 4 / 6 * 100 = 66.7%
	if statsAll.CompletionRate != 66.7 {
		t.Fatalf("all: expected 66.7%% completion rate, got %.1f", statsAll.CompletionRate)
	}
}

func TestStreaks(t *testing.T) {
	now := time.Date(2026, 11, 10, 15, 0, 0, 0, time.Local)

	makeDay := func(year int, month time.Month, day int) time.Time {
		return time.Date(year, month, day, 10, 0, 0, 0, time.Local)
	}

	d10 := makeDay(2026, 11, 10) // Today
	d9 := makeDay(2026, 11, 9)   // Yesterday
	d8 := makeDay(2026, 11, 8)   // 2 days ago
	d7 := makeDay(2026, 11, 7)   // 3 days ago

	d3 := makeDay(2026, 11, 3)
	d2 := makeDay(2026, 11, 2)
	d1 := makeDay(2026, 11, 1)
	d0 := makeDay(2026, 10, 31)
	dm1 := makeDay(2026, 10, 30) // Historical 5-day streak: Oct 30 - Nov 3

	todos := Todos{
		// Historical 5-day streak
		{Title: "H1", Completed: true, CreatedAt: dm1, CompletedAt: &dm1},
		{Title: "H2", Completed: true, CreatedAt: d0, CompletedAt: &d0},
		{Title: "H3", Completed: true, CreatedAt: d1, CompletedAt: &d1},
		{Title: "H4", Completed: true, CreatedAt: d2, CompletedAt: &d2},
		{Title: "H5", Completed: true, CreatedAt: d3, CompletedAt: &d3},

		// Current 4-day streak: Nov 7, 8, 9, 10
		{Title: "C1", Completed: true, CreatedAt: d7, CompletedAt: &d7},
		{Title: "C2", Completed: true, CreatedAt: d8, CompletedAt: &d8},
		{Title: "C3", Completed: true, CreatedAt: d9, CompletedAt: &d9},
		{Title: "C4", Completed: true, CreatedAt: d10, CompletedAt: &d10},
	}

	curr, longest := calculateStreaks(todos, now)
	if curr != 4 {
		t.Fatalf("expected current streak 4, got %d", curr)
	}
	if longest != 5 {
		t.Fatalf("expected longest streak 5, got %d", longest)
	}

	// If today has no completion yet, streak is maintained if yesterday had completion
	todosWithoutToday := todos[:len(todos)-1]
	curr2, longest2 := calculateStreaks(todosWithoutToday, now)
	if curr2 != 3 {
		t.Fatalf("expected maintained current streak 3 (from yesterday), got %d", curr2)
	}
	if longest2 != 5 {
		t.Fatalf("expected longest streak 5, got %d", longest2)
	}

	// If neither today nor yesterday has completion, current streak is 0
	todosBroken := todos[:len(todos)-2] // stops at d8 (2 days ago)
	curr3, longest3 := calculateStreaks(todosBroken, now)
	if curr3 != 0 {
		t.Fatalf("expected broken current streak 0, got %d", curr3)
	}
	if longest3 != 5 {
		t.Fatalf("expected longest streak 5, got %d", longest3)
	}
}

func TestFormatStats(t *testing.T) {
	stats := TodoStats{
		Period:         "week",
		PeriodLabel:    "Last 7 days",
		TasksCompleted: 8,
		TasksCreated:   10,
		CompletionRate: 80.0,
		DailyAverage:   2.5,
		CurrentStreak:  3,
		LongestStreak:  7,
		TagBreakdown: map[string]int{
			"#work": 5,
			"#life": 3,
		},
	}

	output := FormatStats(stats)
	if !strings.Contains(output, "Todo Statistics (Last 7 days)") {
		t.Fatalf("missing header in output: %s", output)
	}
	if !strings.Contains(output, "Tasks Completed:     8") {
		t.Fatalf("missing completed tasks in output: %s", output)
	}
	if !strings.Contains(output, "Tasks Created:       10") {
		t.Fatalf("missing created tasks in output: %s", output)
	}
	if !strings.Contains(output, "Completion Rate:     80.0%") {
		t.Fatalf("missing completion rate in output: %s", output)
	}
	if !strings.Contains(output, "Daily Average:       2.5 tasks/day") {
		t.Fatalf("missing daily average in output: %s", output)
	}
	if !strings.Contains(output, "Current Streak:      3 day(s)") {
		t.Fatalf("missing current streak in output: %s", output)
	}
	if !strings.Contains(output, "Longest Streak:      7 day(s)") {
		t.Fatalf("missing longest streak in output: %s", output)
	}
	if !strings.Contains(output, "Tag Breakdown:") {
		t.Fatalf("missing tag breakdown in output: %s", output)
	}
	if !strings.Contains(output, "#life:") || !strings.Contains(output, "#work:") {
		t.Fatalf("missing specific tags in output: %s", output)
	}
}

func TestOutputStats_JSON(t *testing.T) {
	stats := TodoStats{
		Period:         "all",
		PeriodLabel:    "All time",
		TasksCompleted: 5,
		TasksCreated:   5,
		CompletionRate: 100.0,
		DailyAverage:   1.0,
		CurrentStreak:  2,
		LongestStreak:  2,
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := OutputStats(stats, true)
	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("OutputStats returned error: %v", err)
	}

	var buf bytes.Buffer
	io.Copy(&buf, r)

	var unmarshaled TodoStats
	if err := json.Unmarshal(buf.Bytes(), &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal JSON output %q: %v", buf.String(), err)
	}

	if unmarshaled.Period != "all" || unmarshaled.TasksCompleted != 5 || unmarshaled.CompletionRate != 100.0 {
		t.Fatalf("unmarshaled values mismatch: %+v", unmarshaled)
	}
}

func TestRunCommand_StatsIntegration(t *testing.T) {
	now := time.Now()
	todos := Todos{
		{Title: "Task 1 #cli", Completed: true, CreatedAt: now, CompletedAt: &now},
		{Title: "Task 2 #cli", Completed: false, CreatedAt: now, CompletedAt: nil},
	}

	// Subcommand forms
	subcommandTests := [][]string{
		{"stats"},
		{"stats", "today"},
		{"stats", "week"},
		{"stats", "month"},
		{"stats", "all"},
		{"stats", "-p", "week"},
		{"stats", "--period", "month"},
		{"stats", "--json"},
		{"stats", "week", "--json"},
		{"stats", "--json", "week"},
	}

	for _, cmd := range subcommandTests {
		mutated, err := RunCommand(cmd, &todos)
		if err != nil {
			t.Fatalf("RunCommand(%v) returned error: %v", cmd, err)
		}
		if mutated {
			t.Fatalf("RunCommand(%v) should not mutate todos", cmd)
		}
	}

	// Flag forms
	flagTests := [][]string{
		{"--stats"},
		{"--stats", "-p", "week"},
		{"--stats", "--period", "month", "--json"},
		{"-p", "today"},
	}

	for _, cmd := range flagTests {
		mutated, err := RunCommand(cmd, &todos)
		if err != nil {
			t.Fatalf("RunCommand(%v) returned error: %v", cmd, err)
		}
		if mutated {
			t.Fatalf("RunCommand(%v) should not mutate todos", cmd)
		}
	}

	// Invalid period
	_, err := RunCommand([]string{"stats", "-p", "invalid_period"}, &todos)
	if err == nil {
		t.Fatalf("expected error for invalid period, got nil")
	}
}
