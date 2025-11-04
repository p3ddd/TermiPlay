package models

import (
	"fmt"
	"strings"

	"termiplay/go-backend/game"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	borderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("238"))

	cellStyle = lipgloss.NewStyle().
			Width(3).
			Align(lipgloss.Center).
			Padding(0, 0)

	hiddenStyle = cellStyle.Copy().
			Background(lipgloss.Color("240")).
			Foreground(lipgloss.Color("255"))

	revealedStyle = cellStyle.Copy().
			Background(lipgloss.Color("235")).
			Foreground(lipgloss.Color("255"))

	flagStyle = cellStyle.Copy().
			Background(lipgloss.Color("202")).
			Foreground(lipgloss.Color("255")).
			Bold(true)

	cursorStyle = cellStyle.Copy().
			Background(lipgloss.Color("205")).
			Foreground(lipgloss.Color("255")).
			Bold(true)

	cursorRevealedStyle = cellStyle.Copy().
				Background(lipgloss.Color("33")).
				Foreground(lipgloss.Color("255")).
				Bold(true)

	mineStyle = cellStyle.Copy().
			Background(lipgloss.Color("196")).
			Foreground(lipgloss.Color("255")).
			Bold(true)

	minesweeperInfoStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				MarginTop(1)

	minesweeperHelpStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				MarginTop(1)
)

type MinesweeperModel struct {
	game       *game.Minesweeper
	cursorX    int
	cursorY    int
	difficulty game.Difficulty
	showWin    bool
}

func NewMinesweeperModel(difficulty game.Difficulty) *MinesweeperModel {
	return &MinesweeperModel{
		game:       game.NewMinesweeper(difficulty),
		cursorX:    0,
		cursorY:    0,
		difficulty: difficulty,
		showWin:    false,
	}
}

func (m *MinesweeperModel) Init() tea.Cmd {
	return nil
}

func (m *MinesweeperModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.game.GameOver && !m.showWin {
		m.showWin = true
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k", "w":
			if m.cursorY > 0 {
				m.cursorY--
			}
		case "down", "j", "s":
			if m.cursorY < m.game.Height-1 {
				m.cursorY++
			}
		case "left", "h", "a":
			if m.cursorX > 0 {
				m.cursorX--
			}
		case "right", "l", "d":
			if m.cursorX < m.game.Width-1 {
				m.cursorX++
			}
		case " ", "enter":
			if !m.game.GameOver {
				m.game.Reveal(m.cursorX, m.cursorY)
			}
		case "f":
			if !m.game.GameOver {
				m.game.ToggleFlag(m.cursorX, m.cursorY)
			}
		case "r":
			if m.game.GameOver {
				m.game = game.NewMinesweeper(m.difficulty)
				m.cursorX = 0
				m.cursorY = 0
				m.showWin = false
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *MinesweeperModel) View() string {
	if m.game.GameOver && m.showWin {
		return m.renderGameOver()
	}

	var b strings.Builder

	// 游戏信息
	elapsed := m.game.GetElapsedTime()
	info := fmt.Sprintf("雷数: %d | 标记: %d | 时间: %d秒",
		m.game.MineCount,
		m.game.Flags,
		int(elapsed.Seconds()))

	b.WriteString(minesweeperInfoStyle.Render(info))
	b.WriteString("\n\n")

	// 游戏网格
	grid := make([]string, m.game.Height)
	for y := 0; y < m.game.Height; y++ {
		row := make([]string, m.game.Width)
		for x := 0; x < m.game.Width; x++ {
			cell := m.game.Grid[y][x]
			var cellStr string
			var style lipgloss.Style

			if x == m.cursorX && y == m.cursorY && !m.game.GameOver {
				switch cell.State {
				case game.CellFlagged:
					style = cursorStyle.Copy().Background(lipgloss.Color("202"))
					cellStr = "🚩"
				case game.CellRevealed:
					// 已解开的区域使用特殊的光标样式
					style = cursorRevealedStyle.Copy()
					content := m.getCellContentPlain(cell)
					// 保持数字的颜色
					if cell.Adjacent > 0 && cell.Adjacent <= 8 {
						colors := []lipgloss.Color{
							lipgloss.Color("39"),  // 1 - 蓝色
							lipgloss.Color("46"),  // 2 - 绿色
							lipgloss.Color("196"), // 3 - 红色
							lipgloss.Color("21"),  // 4 - 深蓝
							lipgloss.Color("124"), // 5 - 深红
							lipgloss.Color("45"),  // 6 - 青色
							lipgloss.Color("0"),   // 7 - 黑色
							lipgloss.Color("240"), // 8 - 灰色
						}
						style = style.Foreground(colors[cell.Adjacent-1])
					}
					cellStr = content
				default:
					style = cursorStyle
					cellStr = "?"
				}
			} else {
				switch cell.State {
				case game.CellHidden:
					style = hiddenStyle
					cellStr = " "
				case game.CellFlagged:
					style = flagStyle
					cellStr = "🚩"
				case game.CellRevealed:
					if cell.IsMine {
						style = mineStyle
						cellStr = "💣"
					} else {
						style = revealedStyle
						cellStr = m.getCellContent(cell)
					}
				}
			}

			// 确保所有单元格都使用相同的宽度设置
			style = style.Width(3).Align(lipgloss.Center)

			row[x] = style.Width(3).Align(lipgloss.Center).Render(cellStr)
		}
		grid[y] = strings.Join(row, "")
	}

	b.WriteString(borderStyle.Render(strings.Join(grid, "\n")))
	b.WriteString("\n\n")

	// 帮助信息
	help := "方向键移动 | 空格/Enter 翻开 | F 标记 | R 重玩 | Q 退出"
	b.WriteString(minesweeperHelpStyle.Render(help))

	return b.String()
}

func (m *MinesweeperModel) getCellContent(cell game.Cell) string {
	if cell.IsMine {
		return "💣"
	}
	if cell.Adjacent == 0 {
		return " "
	}
	colors := []lipgloss.Color{
		lipgloss.Color("39"),  // 1 - 蓝色
		lipgloss.Color("46"),  // 2 - 绿色
		lipgloss.Color("196"), // 3 - 红色
		lipgloss.Color("21"),  // 4 - 深蓝
		lipgloss.Color("124"), // 5 - 深红
		lipgloss.Color("45"),  // 6 - 青色
		lipgloss.Color("0"),   // 7 - 黑色
		lipgloss.Color("240"), // 8 - 灰色
	}
	if cell.Adjacent <= 8 {
		style := lipgloss.NewStyle().Foreground(colors[cell.Adjacent-1]).Bold(true)
		return style.Render(fmt.Sprintf("%d", cell.Adjacent))
	}
	return fmt.Sprintf("%d", cell.Adjacent)
}

// getCellContentPlain 返回纯文本内容，不包含样式
func (m *MinesweeperModel) getCellContentPlain(cell game.Cell) string {
	if cell.IsMine {
		return "💣"
	}
	if cell.Adjacent == 0 {
		return " "
	}
	return fmt.Sprintf("%d", cell.Adjacent)
}

func (m *MinesweeperModel) renderGameOver() string {
	var b strings.Builder

	if m.game.Won {
		winMsg := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("46")).
			MarginBottom(1).
			Align(lipgloss.Center).
			Render("🎉 恭喜！你赢了！ 🎉")
		b.WriteString(winMsg)
	} else {
		loseMsg := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("196")).
			MarginBottom(1).
			Align(lipgloss.Center).
			Render("💥 游戏结束！你踩到雷了！ 💥")
		b.WriteString(loseMsg)
	}

	b.WriteString("\n\n")

	// 显示完整网格
	grid := make([]string, m.game.Height)
	for y := 0; y < m.game.Height; y++ {
		row := make([]string, m.game.Width)
		for x := 0; x < m.game.Width; x++ {
			cell := m.game.Grid[y][x]
			var cellStr string
			var style lipgloss.Style

			if cell.IsMine {
				if cell.State == game.CellFlagged {
					style = flagStyle
					cellStr = "🚩"
				} else {
					style = mineStyle
					cellStr = "💣"
				}
			} else {
				switch cell.State {
				case game.CellFlagged:
					style = flagStyle
					cellStr = "🚩"
				default:
					style = revealedStyle
					cellStr = m.getCellContent(cell)
				}
			}

			row[x] = style.Width(3).Align(lipgloss.Center).Render(cellStr)
		}
		grid[y] = strings.Join(row, "")
	}

	b.WriteString(borderStyle.Render(strings.Join(grid, "\n")))
	b.WriteString("\n\n")

	elapsed := m.game.GetElapsedTime()
	stats := fmt.Sprintf("用时: %d秒", int(elapsed.Seconds()))
	b.WriteString(minesweeperInfoStyle.Render(stats))
	b.WriteString("\n\n")

	help := "R 重新开始 | Q 返回大厅"
	b.WriteString(minesweeperHelpStyle.Render(help))

	return b.String()
}
