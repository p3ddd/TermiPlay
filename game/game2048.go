package game

import (
	"math/rand"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

type Game2048 struct {
	Grid     [4][4]int
	Score    int
	GameOver bool
	Won      bool
}

func NewGame2048() *Game2048 {
	g := &Game2048{
		Score:    0,
		GameOver: false,
		Won:      false,
	}
	g.addRandomTile()
	g.addRandomTile()
	return g
}

func (g *Game2048) addRandomTile() {
	var empty []struct{ x, y int }
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			if g.Grid[y][x] == 0 {
				empty = append(empty, struct{ x, y int }{x, y})
			}
		}
	}

	if len(empty) == 0 {
		return
	}

	pos := empty[rng.Intn(len(empty))]
	if rng.Float32() < 0.9 {
		g.Grid[pos.y][pos.x] = 2
	} else {
		g.Grid[pos.y][pos.x] = 4
	}
}

func (g *Game2048) Move(direction string) bool {
	changed := false

	switch direction {
	case "up":
		changed = g.moveUp()
	case "down":
		changed = g.moveDown()
	case "left":
		changed = g.moveLeft()
	case "right":
		changed = g.moveRight()
	}

	if changed {
		g.addRandomTile()
		g.checkGameState()
	}

	return changed
}

func (g *Game2048) moveLeft() bool {
	changed := false
	for y := 0; y < 4; y++ {
		row := make([]int, 0, 4)
		for x := 0; x < 4; x++ {
			if g.Grid[y][x] != 0 {
				row = append(row, g.Grid[y][x])
			}
		}

		merged := make([]int, 0, 4)
		for i := 0; i < len(row); i++ {
			if i < len(row)-1 && row[i] == row[i+1] {
				merged = append(merged, row[i]*2)
				g.Score += row[i] * 2
				if row[i]*2 == 2048 && !g.Won {
					g.Won = true
				}
				i++
				changed = true
			} else {
				merged = append(merged, row[i])
			}
		}

		for x := 0; x < 4; x++ {
			if x < len(merged) {
				if g.Grid[y][x] != merged[x] {
					changed = true
				}
				g.Grid[y][x] = merged[x]
			} else {
				if g.Grid[y][x] != 0 {
					changed = true
				}
				g.Grid[y][x] = 0
			}
		}
	}
	return changed
}

func (g *Game2048) moveRight() bool {
	changed := false
	for y := 0; y < 4; y++ {
		row := make([]int, 0, 4)
		for x := 3; x >= 0; x-- {
			if g.Grid[y][x] != 0 {
				row = append(row, g.Grid[y][x])
			}
		}

		merged := make([]int, 0, 4)
		for i := 0; i < len(row); i++ {
			if i < len(row)-1 && row[i] == row[i+1] {
				merged = append(merged, row[i]*2)
				g.Score += row[i] * 2
				if row[i]*2 == 2048 && !g.Won {
					g.Won = true
				}
				i++
				changed = true
			} else {
				merged = append(merged, row[i])
			}
		}

		for x := 3; x >= 0; x-- {
			idx := 3 - x
			if idx < len(merged) {
				if g.Grid[y][x] != merged[idx] {
					changed = true
				}
				g.Grid[y][x] = merged[idx]
			} else {
				if g.Grid[y][x] != 0 {
					changed = true
				}
				g.Grid[y][x] = 0
			}
		}
	}
	return changed
}

func (g *Game2048) moveUp() bool {
	changed := false
	for x := 0; x < 4; x++ {
		col := make([]int, 0, 4)
		for y := 0; y < 4; y++ {
			if g.Grid[y][x] != 0 {
				col = append(col, g.Grid[y][x])
			}
		}

		merged := make([]int, 0, 4)
		for i := 0; i < len(col); i++ {
			if i < len(col)-1 && col[i] == col[i+1] {
				merged = append(merged, col[i]*2)
				g.Score += col[i] * 2
				if col[i]*2 == 2048 && !g.Won {
					g.Won = true
				}
				i++
				changed = true
			} else {
				merged = append(merged, col[i])
			}
		}

		for y := 0; y < 4; y++ {
			if y < len(merged) {
				if g.Grid[y][x] != merged[y] {
					changed = true
				}
				g.Grid[y][x] = merged[y]
			} else {
				if g.Grid[y][x] != 0 {
					changed = true
				}
				g.Grid[y][x] = 0
			}
		}
	}
	return changed
}

func (g *Game2048) moveDown() bool {
	changed := false
	for x := 0; x < 4; x++ {
		col := make([]int, 0, 4)
		for y := 3; y >= 0; y-- {
			if g.Grid[y][x] != 0 {
				col = append(col, g.Grid[y][x])
			}
		}

		merged := make([]int, 0, 4)
		for i := 0; i < len(col); i++ {
			if i < len(col)-1 && col[i] == col[i+1] {
				merged = append(merged, col[i]*2)
				g.Score += col[i] * 2
				if col[i]*2 == 2048 && !g.Won {
					g.Won = true
				}
				i++
				changed = true
			} else {
				merged = append(merged, col[i])
			}
		}

		for y := 3; y >= 0; y-- {
			idx := 3 - y
			if idx < len(merged) {
				if g.Grid[y][x] != merged[idx] {
					changed = true
				}
				g.Grid[y][x] = merged[idx]
			} else {
				if g.Grid[y][x] != 0 {
					changed = true
				}
				g.Grid[y][x] = 0
			}
		}
	}
	return changed
}

func (g *Game2048) checkGameState() {
	// 检查是否还有空格
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			if g.Grid[y][x] == 0 {
				return
			}
		}
	}

	// 检查是否可以合并
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			current := g.Grid[y][x]
			if y < 3 && g.Grid[y+1][x] == current {
				return
			}
			if x < 3 && g.Grid[y][x+1] == current {
				return
			}
		}
	}

	g.GameOver = true
}

func (g *Game2048) Reset() {
	g.Grid = [4][4]int{}
	g.Score = 0
	g.GameOver = false
	g.Won = false
	g.addRandomTile()
	g.addRandomTile()
}
