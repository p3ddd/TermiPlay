package models

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	Minesweeper = iota
	Game2048
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			MarginBottom(1).
			Align(lipgloss.Center)

	menuItemStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(lipgloss.Color("240"))

	selectedStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(lipgloss.Color("205")).
			Bold(true)

	lobbyHelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(2)
)

type LobbyModel struct {
	choices    []string
	cursor     int
	selected   int
	gameChosen bool
}

func NewLobbyModel() *LobbyModel {
	return &LobbyModel{
		choices: []string{"扫雷 (Minesweeper)", "2048"},
		cursor:  0,
	}
}

func (m *LobbyModel) Init() tea.Cmd {
	return nil
}

func (m *LobbyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			m.selected = m.cursor
			m.gameChosen = true
		case "q", "ctrl+c":
			m.selected = -1
			m.gameChosen = true
		}
	}
	return m, nil
}

func (m *LobbyModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("🎮 欢迎来到 TermiPlay 游戏大厅 🎮"))
	b.WriteString("\n\n")
	b.WriteString("请选择游戏：\n\n")

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		var style lipgloss.Style
		if m.cursor == i {
			style = selectedStyle
		} else {
			style = menuItemStyle
		}

		b.WriteString(fmt.Sprintf("%s %s\n", cursor, style.Render(choice)))
	}

	b.WriteString("\n")
	b.WriteString(lobbyHelpStyle.Render("↑/↓ 选择 | Enter 确认 | q 退出"))

	return b.String()
}

func (m *LobbyModel) GetSelected() int {
	return m.selected
}

func (m *LobbyModel) IsGameChosen() bool {
	return m.gameChosen
}
