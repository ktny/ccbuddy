package tui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ktny/ccbuddy/internal/buddy"
	"github.com/ktny/ccbuddy/internal/storage"
)

// Model represents the TUI application state
type Model struct {
	buddy *buddy.Buddy
	store *storage.Store
	err   error
}

// NewModel creates a new TUI model
func NewModel() *Model {
	return &Model{
		store: storage.NewStore(),
	}
}

// Init is the first function that will be called
func (m *Model) Init() tea.Cmd {
	return m.loadBuddy
}

// Update handles messages and returns the updated model
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "r":
			return m, m.loadBuddy
		case "f":
			if m.buddy != nil {
				m.buddy.Feed()
				return m, m.saveBuddy
			}
		case "h":
			if m.buddy != nil && m.buddy.State == buddy.StateEgg {
				if err := m.buddy.Hatch(); err == nil {
					return m, m.saveBuddy
				}
			}
		}
	case buddyLoadedMsg:
		m.buddy = msg.buddy
		m.err = msg.err
	case buddySavedMsg:
		m.err = msg.err
	}

	return m, nil
}

// View renders the TUI
func (m *Model) View() string {
	var content string

	if m.err != nil {
		content = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Render(fmt.Sprintf("エラー: %v", m.err))
	} else if m.buddy == nil {
		content = "buddyが見つかりません。新しいbuddyを作成するには 'n' を押してください。"
	} else {
		content = m.renderBuddy()
	}

	helpText := "\n\nキー操作:\n" +
		"• q: 終了\n" +
		"• r: リロード"

	if m.buddy != nil {
		helpText += "\n• f: 餌やり"
		if m.buddy.State == buddy.StateEgg {
			helpText += "\n• h: 孵化"
		}
	}

	return content + helpText
}

// renderBuddy renders the buddy's current state
func (m *Model) renderBuddy() string {
	if m.buddy == nil {
		return ""
	}

	// ステータス表示
	status := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		Render("🥚 ccbuddy ステータス")

	// 状態表示
	var stateIcon string
	switch m.buddy.State {
	case buddy.StateEgg:
		stateIcon = "🥚"
	case buddy.StateHatched:
		stateIcon = "🐣"
	default:
		stateIcon = "❓"
	}

	stateText := fmt.Sprintf("%s 状態: %s", stateIcon, m.buddy.State)

	// 健康状態表示（プログレスバー風）
	healthBar := m.renderHealthBar(m.buddy.Health)
	healthText := fmt.Sprintf("❤️  健康状態: %d/100", m.buddy.Health)

	// 年齢表示
	age := m.buddy.Age()
	ageText := fmt.Sprintf("⏰ 年齢: %s", formatDuration(age))

	// 最後の餌やり時刻
	timeSinceFed := time.Since(m.buddy.LastFedAt)
	fedText := fmt.Sprintf("🍽️  最後の餌やり: %s前", formatDuration(timeSinceFed))

	// レイアウト
	content := fmt.Sprintf(
		"%s\n\n%s\n%s\n%s\n\n%s\n%s",
		status,
		stateText,
		healthText,
		healthBar,
		ageText,
		fedText,
	)

	return content
}

// renderHealthBar creates a visual health bar
func (m *Model) renderHealthBar(health int) string {
	barWidth := 20
	filledWidth := (health * barWidth) / 100

	var style lipgloss.Style
	if health > 70 {
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("10")) // green
	} else if health > 30 {
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("11")) // yellow
	} else {
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("9")) // red
	}

	filled := style.Render(string(make([]rune, filledWidth)))
	empty := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Render(string(make([]rune, barWidth-filledWidth)))

	return fmt.Sprintf("[%s%s]", filled, empty)
}

// formatDuration formats a duration in a human-readable way
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0f秒", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.0f分", d.Minutes())
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%.1f時間", d.Hours())
	}
	return fmt.Sprintf("%.1f日", d.Hours()/24)
}

// Messages for async operations
type buddyLoadedMsg struct {
	buddy *buddy.Buddy
	err   error
}

type buddySavedMsg struct {
	err error
}

// loadBuddy loads the buddy from storage
func (m *Model) loadBuddy() tea.Msg {
	if m.store.BuddyExists() {
		buddy, err := m.store.LoadBuddy()
		return buddyLoadedMsg{buddy: buddy, err: err}
	}
	return buddyLoadedMsg{buddy: nil, err: nil}
}

// saveBuddy saves the buddy to storage
func (m *Model) saveBuddy() tea.Msg {
	if m.buddy != nil {
		err := m.store.SaveBuddy(m.buddy)
		return buddySavedMsg{err: err}
	}
	return buddySavedMsg{err: nil}
}
