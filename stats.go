package main

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"
	"time"
)

type TodoStats struct {
	Period         string         `json:"period"`
	PeriodLabel    string         `json:"period_label"`
	From           *time.Time     `json:"from,omitempty"`
	To             *time.Time     `json:"to,omitempty"`
	TasksCompleted int            `json:"tasks_completed"`
	TasksCreated   int            `json:"tasks_created"`
	CompletionRate float64        `json:"completion_rate"`
	DailyAverage   float64        `json:"daily_average"`
	CurrentStreak  int            `json:"current_streak"`
	LongestStreak  int            `json:"longest_streak"`
	TagBreakdown   map[string]int `json:"tag_breakdown,omitempty"`
}

var tagRegex = regexp.MustCompile(`#([a-zA-Z0-9_\-]+)`)

// periodRange implements the requirement:
// periodRange(period string) (from, to time.Time)
func periodRange(period string) (from, to time.Time) {
	from, to, _, _ = periodRangeAt(period, time.Now())
	return from, to
}

// periodRangeAt computes the date boundaries and display label relative to a reference time.
func periodRangeAt(period string, now time.Time) (from, to time.Time, label string, err error) {
	norm := strings.ToLower(strings.TrimSpace(period))
	if norm == "" {
		norm = "all"
	}

	loc := now.Location()
	year, month, day := now.Date()

	switch norm {
	case "today", "day", "d":
		from = time.Date(year, month, day, 0, 0, 0, 0, loc)
		to = now
		label = "Today"
	case "week", "w", "7d":
		from = now.AddDate(0, 0, -7)
		to = now
		label = "Last 7 days"
	case "month", "m", "30d":
		from = now.AddDate(0, 0, -30)
		to = now
		label = "Last 30 days"
	case "all", "alltime", "all-time":
		from = time.Time{}
		to = now
		label = "All time"
	default:
		return time.Time{}, time.Time{}, "", fmt.Errorf("invalid period %q: must be one of 'today', 'week', 'month', or 'all'", period)
	}

	return from, to, label, nil
}

// ComputeStats calculates complete todo statistics for a given period.
func ComputeStats(todos Todos, period string, now time.Time) (TodoStats, error) {
	from, to, label, err := periodRangeAt(period, now)
	if err != nil {
		return TodoStats{}, err
	}

	normPeriod := strings.ToLower(strings.TrimSpace(period))
	if normPeriod == "" {
		normPeriod = "all"
	}

	isAll := normPeriod == "all" || normPeriod == "alltime" || normPeriod == "all-time"

	var completedCount int
	var createdCount int
	activeDaysSet := make(map[string]struct{})
	tagCounts := make(map[string]int)

	// Filter tasks for the selected period
	for _, item := range todos {
		// Created count
		if isAll {
			createdCount++
		} else if !item.CreatedAt.IsZero() && !item.CreatedAt.Before(from) && !item.CreatedAt.After(to) {
			createdCount++
		}

		// Completed count
		if item.Completed {
			var isCompletedInPeriod bool
			if isAll {
				isCompletedInPeriod = true
			} else if item.CompletedAt != nil && !item.CompletedAt.Before(from) && !item.CompletedAt.After(to) {
				isCompletedInPeriod = true
			}

			if isCompletedInPeriod {
				completedCount++

				// Active day tracking
				if item.CompletedAt != nil {
					dayKey := item.CompletedAt.In(now.Location()).Format("2006-01-02")
					activeDaysSet[dayKey] = struct{}{}
				}

				// Extract tags from completed task title
				matches := tagRegex.FindAllStringSubmatch(item.Title, -1)
				for _, m := range matches {
					if len(m) > 1 {
						tagName := strings.ToLower(m[1])
						tagCounts["#"+tagName]++
					}
				}
			}
		}
	}

	// Completion Rate: (Completed / Created) * 100%
	var rate float64
	if createdCount > 0 {
		rate = (float64(completedCount) / float64(createdCount)) * 100.0
		rate = math.Round(rate*10) / 10
	} else if completedCount > 0 {
		rate = 100.0
	}

	// Daily Average: Average tasks completed per active day
	var dailyAvg float64
	activeDays := len(activeDaysSet)
	if activeDays > 0 {
		dailyAvg = float64(completedCount) / float64(activeDays)
		dailyAvg = math.Round(dailyAvg*10) / 10
	}

	// Streaks calculation across all history
	currStreak, longestStreak := calculateStreaks(todos, now)

	var fromPtr, toPtr *time.Time
	if !isAll {
		fromPtr = &from
		toPtr = &to
	}

	var tagsMap map[string]int
	if len(tagCounts) > 0 {
		tagsMap = tagCounts
	}

	return TodoStats{
		Period:         normPeriod,
		PeriodLabel:    label,
		From:           fromPtr,
		To:             toPtr,
		TasksCompleted: completedCount,
		TasksCreated:   createdCount,
		CompletionRate: rate,
		DailyAverage:   dailyAvg,
		CurrentStreak:  currStreak,
		LongestStreak:  longestStreak,
		TagBreakdown:   tagsMap,
	}, nil
}

// calculateStreaks returns the current streak and longest historical streak.
func calculateStreaks(todos Todos, now time.Time) (currentStreak int, longestStreak int) {
	loc := now.Location()
	completionDates := make(map[string]time.Time)

	for _, item := range todos {
		if item.Completed && item.CompletedAt != nil {
			local := item.CompletedAt.In(loc)
			normalized := time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, loc)
			dateKey := normalized.Format("2006-01-02")
			completionDates[dateKey] = normalized
		}
	}

	if len(completionDates) == 0 {
		return 0, 0
	}

	// Sort unique completion days chronologically
	var days []time.Time
	for _, d := range completionDates {
		days = append(days, d)
	}
	sort.Slice(days, func(i, j int) bool {
		return days[i].Before(days[j])
	})

	// Longest streak
	longest := 1
	currentRun := 1
	for i := 1; i < len(days); i++ {
		prevPlusOne := days[i-1].AddDate(0, 0, 1)
		if days[i].Equal(prevPlusOne) {
			currentRun++
			if currentRun > longest {
				longest = currentRun
			}
		} else {
			currentRun = 1
		}
	}

	// Current streak: consecutive days leading up to today
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	todayKey := today.Format("2006-01-02")
	yesterday := today.AddDate(0, 0, -1)
	yesterdayKey := yesterday.Format("2006-01-02")

	current := 0
	checkDay := today

	// If today has completion, streak starts from today.
	// If today does not yet have completion, streak is maintained if yesterday had completion.
	if _, ok := completionDates[todayKey]; ok {
		checkDay = today
	} else if _, ok := completionDates[yesterdayKey]; ok {
		checkDay = yesterday
	} else {
		// Streak is broken
		return 0, longest
	}

	for {
		key := checkDay.Format("2006-01-02")
		if _, ok := completionDates[key]; ok {
			current++
			checkDay = checkDay.AddDate(0, 0, -1)
		} else {
			break
		}
	}

	if current > longest {
		longest = current
	}

	return current, longest
}

// FormatStats produces aligned, clean terminal text output with dividers.
func FormatStats(s TodoStats) string {
	var b strings.Builder

	header := fmt.Sprintf("Todo Statistics (%s)", s.PeriodLabel)
	b.WriteString(header + "\n")
	b.WriteString(strings.Repeat("─", 40) + "\n")

	b.WriteString(fmt.Sprintf("%-21s%d\n", "Tasks Completed:", s.TasksCompleted))
	b.WriteString(fmt.Sprintf("%-21s%d\n", "Tasks Created:", s.TasksCreated))
	b.WriteString(fmt.Sprintf("%-21s%.1f%%\n", "Completion Rate:", s.CompletionRate))
	b.WriteString(fmt.Sprintf("%-21s%.1f tasks/day\n", "Daily Average:", s.DailyAverage))
	b.WriteString(fmt.Sprintf("%-21s%d day(s)\n", "Current Streak:", s.CurrentStreak))
	b.WriteString(fmt.Sprintf("%-21s%d day(s)", "Longest Streak:", s.LongestStreak))

	if len(s.TagBreakdown) > 0 {
		b.WriteString("\n" + strings.Repeat("─", 40) + "\n")
		b.WriteString("Tag Breakdown:\n")

		// Sort tags alphabetically
		var tags []string
		for tag := range s.TagBreakdown {
			tags = append(tags, tag)
		}
		sort.Strings(tags)

		for _, tag := range tags {
			b.WriteString(fmt.Sprintf("  %-19s%d completed\n", tag+":", s.TagBreakdown[tag]))
		}
		// Trim trailing newline for consistency
		return strings.TrimRight(b.String(), "\n")
	}

	return b.String()
}

// OutputStats formats and prints the stats to stdout either as JSON or as formatted text.
func OutputStats(s TodoStats, jsonOutput bool) error {
	if jsonOutput {
		data, err := json.MarshalIndent(s, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}

	fmt.Println(FormatStats(s))
	return nil
}
