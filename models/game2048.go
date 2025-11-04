package models

import (
	"fmt"
	"strings"

	"termiplay/go-backend/game"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	gridStyle = lipgloss.NewStyle().
			Padding(1, 1)

	game2048CellStyle = lipgloss.NewStyle().
				Width(10).
				Height(3).
				Align(lipgloss.Center).
				AlignVertical(lipgloss.Center).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("238"))

	game2048InfoStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				MarginBottom(1)

	game2048HelpStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				MarginTop(1)

	winStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("46")).
			MarginBottom(1).
			Align(lipgloss.Center)

	gameOverStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("196")).
			MarginBottom(1).
			Align(lipgloss.Center)
)

func getCellColor(value int) lipgloss.Color {
	colors := map[int]lipgloss.Color{
		0:    lipgloss.Color("235"),
		2:    lipgloss.Color("237"),
		4:    lipgloss.Color("238"),
		8:    lipgloss.Color("239"),
		16:   lipgloss.Color("240"),
		32:   lipgloss.Color("241"),
		64:   lipgloss.Color("202"),
		128:  lipgloss.Color("214"),
		256:  lipgloss.Color("226"),
		512:  lipgloss.Color("220"),
		1024: lipgloss.Color("11"),
		2048: lipgloss.Color("196"),
	}
	if color, ok := colors[value]; ok {
		return color
	}
	return lipgloss.Color("235")
}

func getTextColor(value int) lipgloss.Color {
	if value <= 4 {
		return lipgloss.Color("255")
	}
	return lipgloss.Color("0")
}

type Game2048Model struct {
	game *game.Game2048
}

func NewGame2048Model() *Game2048Model {
	return &Game2048Model{
		game: game.NewGame2048(),
	}
}

func (m *Game2048Model) Init() tea.Cmd {
	return nil
}

func (m *Game2048Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.game.GameOver {
			switch msg.String() {
			case "r":
				m.game.Reset()
			case "q", "ctrl+c":
				return m, tea.Quit
			}
		} else {
			switch msg.String() {
			case "up", "k", "w":
				m.game.Move("up")
			case "down", "j", "s":
				m.game.Move("down")
			case "left", "h", "a":
				m.game.Move("left")
			case "right", "l", "d":
				m.game.Move("right")
			case "r":
				m.game.Reset()
			case "q", "ctrl+c":
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m *Game2048Model) View() string {
	var b strings.Builder

	// 游戏信息
	info := fmt.Sprintf("分数: %d", m.game.Score)
	if m.game.Won && !m.game.GameOver {
		info += " | 🎉 达成2048！"
	}
	b.WriteString(game2048InfoStyle.Render(info))
	b.WriteString("\n\n")

	if m.game.GameOver {
		if m.game.Won {
			b.WriteString(winStyle.Render("🎉 恭喜！你达成了2048！ 🎉"))
		} else {
			b.WriteString(gameOverStyle.Render("游戏结束！无法继续移动"))
		}
		b.WriteString("\n\n")
	}

	// 游戏网格 - 使用lipgloss的JoinHorizontal来确保正确的布局
	gridRows := make([]string, 4)
	for y := 0; y < 4; y++ {
		cells := make([]string, 4)
		for x := 0; x < 4; x++ {
			value := m.game.Grid[y][x]
			cellStr := " "
			if value != 0 {
				cellStr = fmt.Sprintf("%d", value)
			}

			style := game2048CellStyle.Copy().
				Background(getCellColor(value)).
				Foreground(getTextColor(value))

			cells[x] = style.Render(cellStr)
		}
		gridRows[y] = lipgloss.JoinHorizontal(lipgloss.Left, cells...)
	}

	b.WriteString(gridStyle.Render(lipgloss.JoinVertical(lipgloss.Top, gridRows...)))
	b.WriteString("\n\n")

	// 帮助信息
	help := "方向键移动 | R 重新开始 | Q 返回大厅"
	b.WriteString(game2048HelpStyle.Render(help))

	return b.String()
}
