package tada

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FilterState int

const (
	FilterAll FilterState = iota
	FilterActive
	FilterCompleted
)

type tuiModel struct {
	todos           *Todos
	storage         *Storage[Todos]
	filter          FilterState
	filteredIndices []int
	cursor          int
	startIndex      int
	adding          bool
	input           textinput.Model
	width           int
	height          int
	quitting        bool
}

// Sophisticated, cohesive color palette (inspired by Tokyo Night / Catppuccin)
var (
	cPrimary    = lipgloss.Color("#7AA2F7") // Bright blue/indigo
	cAccent     = lipgloss.Color("#BB9AF7") // Soft violet
	cGreen      = lipgloss.Color("#9ECE6A") // Soft emerald
	cPending    = lipgloss.Color("#7DCFFF") // Cyan for pending
	cText       = lipgloss.Color("#C0CAF5") // Crisp light text
	cTextMuted  = lipgloss.Color("#565F89") // Muted gray
	cTextDim    = lipgloss.Color("#414868") // Subtle dark gray
	cBorder     = lipgloss.Color("#3B4261") // Clean border gray
	cBorderFoc  = lipgloss.Color("#7AA2F7") // Focused border
	cRowHover   = lipgloss.Color("#1F2335") // Row highlight bg
	cTabBgAct   = lipgloss.Color("#283457") // Active tab background

	// Container & Card
	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(cBorder).
			Padding(1, 3)

	// Header Styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(cPrimary)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(cTextMuted)

	// Filter Tabs
	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(cPrimary).
			Background(cTabBgAct).
			Padding(0, 2)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(cTextMuted).
				Padding(0, 2)

	// Row & Item Styles
	selectedRowStyle = lipgloss.NewStyle().
				Background(cRowHover)

	cursorActive = lipgloss.NewStyle().
			Bold(true).
			Foreground(cPrimary).
			Render("▸ ")

	cursorEmpty = "  "

	checkboxPending = lipgloss.NewStyle().
			Bold(true).
			Foreground(cPending).
			Render("○ ")

	checkboxCompleted = lipgloss.NewStyle().
				Bold(true).
				Foreground(cGreen).
				Render("● ")

	taskPendingStyle = lipgloss.NewStyle().
				Foreground(cText)

	taskSelectedStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFFFFF"))

	taskCompletedStyle = lipgloss.NewStyle().
				Strikethrough(true).
				Foreground(cTextMuted)

	dateStyle = lipgloss.NewStyle().
			Foreground(cTextDim)

	dateSelectedStyle = lipgloss.NewStyle().
				Foreground(cTextMuted)

	// Input Styles
	inputBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(cBorderFoc).
			Padding(0, 1)

	inputPromptStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(cPrimary)

	// Footer & Helpers
	dividerStyle = lipgloss.NewStyle().
			Foreground(cBorder)

	keyBadgeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(cAccent)

	keyDescStyle = lipgloss.NewStyle().
			Foreground(cTextMuted)

	statsNumStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(cPrimary)

	emptyStateStyle = lipgloss.NewStyle().
			Foreground(cTextMuted).
			Padding(2, 0)
)

func newTUIModel(todos *Todos, storage *Storage[Todos]) tuiModel {
	ti := textinput.New()
	ti.Placeholder = "Write a task name and press Enter..."
	ti.Prompt = ""
	ti.CharLimit = 120
	ti.Width = 46

	m := tuiModel{
		todos:   todos,
		storage: storage,
		filter:  FilterAll,
		input:   ti,
	}
	m.rebuildFiltered()
	return m
}

func (m *tuiModel) rebuildFiltered() {
	m.filteredIndices = nil
	for i, t := range *m.todos {
		switch m.filter {
		case FilterAll:
			m.filteredIndices = append(m.filteredIndices, i)
		case FilterActive:
			if !t.Completed {
				m.filteredIndices = append(m.filteredIndices, i)
			}
		case FilterCompleted:
			if t.Completed {
				m.filteredIndices = append(m.filteredIndices, i)
			}
		}
	}

	if len(m.filteredIndices) == 0 {
		m.cursor = 0
		m.startIndex = 0
	} else {
		if m.cursor >= len(m.filteredIndices) {
			m.cursor = len(m.filteredIndices) - 1
		}
		if m.cursor < 0 {
			m.cursor = 0
		}
	}
}

func (m *tuiModel) countStats() (active, completed, total int) {
	total = len(*m.todos)
	for _, t := range *m.todos {
		if t.Completed {
			completed++
		} else {
			active++
		}
	}
	return active, completed, total
}

func (m tuiModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if m.adding {
			switch msg.String() {
			case "enter":
				title := strings.TrimSpace(m.input.Value())
				if title != "" {
					m.todos.add(title)
					if m.storage != nil {
						_ = m.storage.Save(*m.todos)
					}
					m.rebuildFiltered()
					if m.filter != FilterCompleted {
						m.cursor = len(m.filteredIndices) - 1
					}
				}
				m.adding = false
				m.input.Blur()
				m.input.Reset()
				return m, nil

			case "esc":
				m.adding = false
				m.input.Blur()
				m.input.Reset()
				return m, nil

			default:
				var cmd tea.Cmd
				m.input, cmd = m.input.Update(msg)
				return m, cmd
			}
		}

		// Main navigation and controls
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.filteredIndices)-1 {
				m.cursor++
			}

		case " ", "space", "enter":
			if len(m.filteredIndices) > 0 && m.cursor < len(m.filteredIndices) {
				actualIdx := m.filteredIndices[m.cursor]
				_ = m.todos.toggle(actualIdx)
				if m.storage != nil {
					_ = m.storage.Save(*m.todos)
				}
				m.rebuildFiltered()
			}

		case "a", "n":
			m.adding = true
			m.input.Reset()
			m.input.Focus()
			return m, textinput.Blink

		case "d", "backspace", "x":
			if len(m.filteredIndices) > 0 && m.cursor < len(m.filteredIndices) {
				actualIdx := m.filteredIndices[m.cursor]
				_ = m.todos.delete(actualIdx)
				if m.storage != nil {
					_ = m.storage.Save(*m.todos)
				}
				m.rebuildFiltered()
			}

		case "tab":
			m.filter = (m.filter + 1) % 3
			m.rebuildFiltered()

		case "shift+tab":
			m.filter = (m.filter + 2) % 3
			m.rebuildFiltered()
		}
	}

	return m, nil
}

func (m tuiModel) View() string {
	if m.quitting {
		return ""
	}

	// Dynamic, breathing width: 72 columns default, gracefully shrinking on small windows
	targetWidth := 74
	if m.width > 0 && m.width-4 < targetWidth {
		targetWidth = m.width - 4
		if targetWidth < 40 {
			targetWidth = 40
		}
	}
	innerWidth := targetWidth - 6 // accounting for card padding and borders

	// 1. Header (Title on left, progress stats on right)
	activeCount, completedCount, totalCount := m.countStats()
	titleLeft := lipgloss.JoinHorizontal(
		lipgloss.Center,
		titleStyle.Render("📝 TADA TODO"),
		"  ",
		subtitleStyle.Render("• task manager"),
	)

	var percent int
	if totalCount > 0 {
		percent = (completedCount * 100) / totalCount
	}
	progressRight := statsNumStyle.Render(fmt.Sprintf("%d/%d done", completedCount, totalCount)) +
		subtitleStyle.Render(fmt.Sprintf(" (%d%%)", percent))

	headerGap := innerWidth - lipgloss.Width(titleLeft) - lipgloss.Width(progressRight)
	if headerGap < 2 {
		headerGap = 2
	}
	headerRow := titleLeft + strings.Repeat(" ", headerGap) + progressRight

	// 2. Filter Navigation Bar
	allTab := inactiveTabStyle.Render(fmt.Sprintf("All (%d)", totalCount))
	if m.filter == FilterAll {
		allTab = activeTabStyle.Render(fmt.Sprintf("All (%d)", totalCount))
	}

	activeTab := inactiveTabStyle.Render(fmt.Sprintf("Active (%d)", activeCount))
	if m.filter == FilterActive {
		activeTab = activeTabStyle.Render(fmt.Sprintf("Active (%d)", activeCount))
	}

	completedTab := inactiveTabStyle.Render(fmt.Sprintf("Completed (%d)", completedCount))
	if m.filter == FilterCompleted {
		completedTab = activeTabStyle.Render(fmt.Sprintf("Completed (%d)", completedCount))
	}

	tabBar := lipgloss.JoinHorizontal(lipgloss.Top, allTab, " ", activeTab, " ", completedTab)

	// Clean subtle divider line
	divider := dividerStyle.Render(strings.Repeat("─", innerWidth))

	// 3. Task List Section with scrolling
	maxVisible := 8
	if m.height > 22 {
		maxVisible = m.height - 16
		if maxVisible > 16 {
			maxVisible = 16
		}
	}

	if m.cursor < m.startIndex {
		m.startIndex = m.cursor
	} else if m.cursor >= m.startIndex+maxVisible {
		m.startIndex = m.cursor - maxVisible + 1
	}
	if m.startIndex < 0 {
		m.startIndex = 0
	}

	var taskRows []string
	if len(m.filteredIndices) == 0 {
		var msg string
		switch m.filter {
		case FilterActive:
			msg = "✨ No active tasks! You are all caught up."
		case FilterCompleted:
			msg = "No completed tasks yet. Finish one with Space."
		default:
			msg = "No tasks yet! Press 'a' to add your first task."
		}
		taskRows = append(taskRows, emptyStateStyle.Render("   "+msg))
	} else {
		endIndex := m.startIndex + maxVisible
		if endIndex > len(m.filteredIndices) {
			endIndex = len(m.filteredIndices)
		}

		if m.startIndex > 0 {
			taskRows = append(taskRows, dateStyle.Render(fmt.Sprintf("   ↑ %d more above...", m.startIndex)))
		}

		for i := m.startIndex; i < endIndex; i++ {
			actualIdx := m.filteredIndices[i]
			item := (*m.todos)[actualIdx]
			isSelected := i == m.cursor

			// Indicator & Checkbox
			indicator := cursorEmpty
			if isSelected {
				indicator = cursorActive
			}

			checkbox := checkboxPending
			if item.Completed {
				checkbox = checkboxCompleted
			}

			// Date string (right-aligned)
			dateStr := item.CreatedAt.Format("01/02 15:04")
			var renderedDate string
			if isSelected {
				renderedDate = dateSelectedStyle.Render(dateStr)
			} else {
				renderedDate = dateStyle.Render(dateStr)
			}

			// Exact visual width calculation using lipgloss.Width to guarantee flawless alignment
			fixedPrefix := indicator + checkbox
			prefixWidth := lipgloss.Width(fixedPrefix)
			dateWidth := lipgloss.Width(renderedDate)
			maxTitleWidth := innerWidth - prefixWidth - dateWidth - 4
			if maxTitleWidth < 12 {
				maxTitleWidth = 12
			}

			titleText := item.Title
			if len(titleText) > maxTitleWidth {
				titleText = titleText[:maxTitleWidth-1] + "…"
			}

			var renderedTitle string
			if item.Completed {
				renderedTitle = taskCompletedStyle.Render(titleText)
			} else if isSelected {
				renderedTitle = taskSelectedStyle.Render(titleText)
			} else {
				renderedTitle = taskPendingStyle.Render(titleText)
			}

			titleVisualWidth := lipgloss.Width(renderedTitle)
			spacing := innerWidth - prefixWidth - titleVisualWidth - dateWidth
			if spacing < 2 {
				spacing = 2
			}

			rowContent := fixedPrefix + renderedTitle + strings.Repeat(" ", spacing) + renderedDate

			if isSelected {
				rowContent = selectedRowStyle.Width(innerWidth).Render(rowContent)
			}

			taskRows = append(taskRows, rowContent)
		}

		if endIndex < len(m.filteredIndices) {
			taskRows = append(taskRows, dateStyle.Render(fmt.Sprintf("   ↓ %d more below...", len(m.filteredIndices)-endIndex)))
		}
	}

	taskListView := strings.Join(taskRows, "\n")

	// 4. Inline Add Task Section
	var inputSection string
	if m.adding {
		inputField := inputBoxStyle.Width(innerWidth - 16).Render(m.input.View())
		inputPrompt := inputPromptStyle.Render("+ New Task: ")
		inputRow := lipgloss.JoinHorizontal(lipgloss.Center, inputPrompt, inputField)
		inputHint := keyDescStyle.Render("  press Enter to save • Esc to cancel")
		inputSection = lipgloss.JoinVertical(lipgloss.Left, inputRow, inputHint)
	}

	// 5. Clean, Spaced Footer with Keybindings
	footerHelp := lipgloss.JoinHorizontal(
		lipgloss.Top,
		keyBadgeStyle.Render("↑/↓ "), keyDescStyle.Render("move  "),
		keyBadgeStyle.Render("space "), keyDescStyle.Render("toggle  "),
		keyBadgeStyle.Render("a "), keyDescStyle.Render("add  "),
		keyBadgeStyle.Render("d "), keyDescStyle.Render("delete  "),
		keyBadgeStyle.Render("tab "), keyDescStyle.Render("filter  "),
		keyBadgeStyle.Render("q "), keyDescStyle.Render("quit"),
	)

	// Assemble card components with consistent, airy vertical spacing
	var contentElements []string
	contentElements = append(contentElements, headerRow)
	contentElements = append(contentElements, "")
	contentElements = append(contentElements, tabBar)
	contentElements = append(contentElements, divider)
	contentElements = append(contentElements, taskListView)
	contentElements = append(contentElements, divider)

	if m.adding {
		contentElements = append(contentElements, inputSection)
		contentElements = append(contentElements, "")
	}

	contentElements = append(contentElements, footerHelp)

	cardContent := strings.Join(contentElements, "\n")
	card := cardStyle.Width(targetWidth).Render(cardContent)

	if m.width > 0 && m.height > 0 {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, card)
	}

	return card
}

func RunTUI(todos *Todos, storage *Storage[Todos]) error {
	m := newTUIModel(todos, storage)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

type TUIModel = tuiModel

func NewTUIModel(todos *Todos, storage *Storage[Todos]) TUIModel {
	return newTUIModel(todos, storage)
}

func (m *tuiModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

func (m tuiModel) Cursor() int {
	return m.cursor
}

func (m tuiModel) Filter() FilterState {
	return m.filter
}

func (m tuiModel) FilteredIndices() []int {
	return m.filteredIndices
}

func (m tuiModel) Adding() bool {
	return m.adding
}
