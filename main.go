package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	leftPaneWidth = 20
)

var (
	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("120")).
			Bold(true)

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("99")).
			Bold(true)
)

type model struct {
	icons    []string
	cursor   int
	selected map[int]bool
	width    int
	height   int
	ready    bool
}

func loadIcons() ([]string, error) {
	cachePath := filepath.Join(os.Getenv("HOME"), ".cache/nerdfonts/icons-list.txt")
	file, err := os.Open(cachePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var icons []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		icon := strings.TrimSpace(scanner.Text())
		if icon != "" {
			icons = append(icons, icon)
		}
	}
	return icons, scanner.Err()
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			var selected []string
			for i := range m.icons {
				if m.selected[i] {
					selected = append(selected, m.icons[i])
				}
			}
			if len(selected) > 0 {
				copyToClipboard(strings.Join(selected, " "))
			}
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.icons)-1 {
				m.cursor++
			}

		case "g":
			m.cursor = 0

		case "G":
			m.cursor = len(m.icons) - 1

		case " ", "tab":
			m.selected[m.cursor] = !m.selected[m.cursor]
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
	}

	return m, nil
}

func (m model) View() string {
	if !m.ready {
		return "\n  Loading icons..."
	}

	// Calculate available lines (subtract border + padding)
	availableLines := m.height - 4

	// Left pane content
	leftHeader := headerStyle.Render("Selected")
	var selectedIcons []string
	for i, icon := range m.icons {
		if m.selected[i] {
			selectedIcons = append(selectedIcons, icon)
		}
	}

	var leftBody string
	if len(selectedIcons) > 0 {
		leftBody = selectedStyle.Render(strings.Join(selectedIcons, "  "))
	} else {
		leftBody = dimStyle.Render("(none)")
	}

	// Controls at bottom
	controls := "\n\n" + dimStyle.Render("jk space ↵ q")

	leftContent := leftHeader + "\n\n" + leftBody + controls
	leftPane := borderStyle.Width(leftPaneWidth).Render(leftContent)

	// Right pane content
	rightHeader := headerStyle.Render("Icons")

	// How many icon lines can we show?
	headerLines := 3 // header + newlines
	iconLines := availableLines - headerLines
	if iconLines < 1 {
		iconLines = 1
	}

	// Center cursor in viewport
	start := m.cursor - (iconLines / 2)
	if start < 0 {
		start = 0
	}
	end := start + iconLines
	if end > len(m.icons) {
		end = len(m.icons)
		start = end - iconLines
		if start < 0 {
			start = 0
		}
	}

	// Build icon list
	var lines []string
	for i := start; i < end; i++ {
		icon := m.icons[i]

		// Get hex codepoint
		runes := []rune(icon)
		hex := fmt.Sprintf("U+%04X", runes[0])

		var line string
		if i == m.cursor {
			if m.selected[i] {
				line = selectedStyle.Render("✓ " + icon + "  ") + dimStyle.Render(hex)
			} else {
				line = "→ " + dimStyle.Render(icon + "  " + hex)
			}
		} else {
			if m.selected[i] {
				line = selectedStyle.Render("✓ " + icon + "  ") + dimStyle.Render(hex)
			} else {
				line = "  " + dimStyle.Render(icon + "  " + hex)
			}
		}
		lines = append(lines, line)
	}

	rightContent := rightHeader + "\n\n" + strings.Join(lines, "\n")
	rightPane := borderStyle.Render(rightContent)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftPane, "  ", rightPane)
}

func copyToClipboard(text string) {
	cmd := exec.Command("wl-copy")
	cmd.Stdin = strings.NewReader(text)
	cmd.Run()
}

func main() {
	icons, err := loadIcons()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading icons: %v\n", err)
		os.Exit(1)
	}

	m := model{
		icons:    icons,
		cursor:   0,
		selected: make(map[int]bool),
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
