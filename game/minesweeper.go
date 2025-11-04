package game

import (
	"math/rand"
	"time"
)

type CellState int

const (
	CellHidden CellState = iota
	CellRevealed
	CellFlagged
)

type Cell struct {
	IsMine   bool
	State    CellState
	Adjacent int // 周围雷的数量
}

type Minesweeper struct {
	Grid      [][]Cell
	Width     int
	Height    int
	MineCount int
	Flags     int
	Revealed  int
	GameOver  bool
	Won       bool
	StartTime time.Time
}

type Difficulty int

const (
	Easy Difficulty = iota
	Medium
	Hard
)

func NewMinesweeper(difficulty Difficulty) *Minesweeper {
	var width, height, mineCount int

	switch difficulty {
	case Easy:
		width, height, mineCount = 9, 9, 10
	case Medium:
		width, height, mineCount = 16, 16, 40
	case Hard:
		width, height, mineCount = 30, 16, 99
	default:
		width, height, mineCount = 9, 9, 10
	}

	ms := &Minesweeper{
		Width:     width,
		Height:    height,
		MineCount: mineCount,
		Flags:     0,
		Revealed:  0,
		GameOver:  false,
		Won:       false,
		StartTime: time.Now(),
	}

	ms.Grid = make([][]Cell, height)
	for y := range ms.Grid {
		ms.Grid[y] = make([]Cell, width)
	}

	ms.placeMines()
	ms.calculateAdjacent()

	return ms
}

func (ms *Minesweeper) placeMines() {
	rand.Seed(time.Now().UnixNano())
	minesPlaced := 0

	for minesPlaced < ms.MineCount {
		x := rand.Intn(ms.Width)
		y := rand.Intn(ms.Height)

		if !ms.Grid[y][x].IsMine {
			ms.Grid[y][x].IsMine = true
			minesPlaced++
		}
	}
}

func (ms *Minesweeper) calculateAdjacent() {
	for y := 0; y < ms.Height; y++ {
		for x := 0; x < ms.Width; x++ {
			if ms.Grid[y][x].IsMine {
				continue
			}

			count := 0
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					if dx == 0 && dy == 0 {
						continue
					}
					nx, ny := x+dx, y+dy
					if nx >= 0 && nx < ms.Width && ny >= 0 && ny < ms.Height {
						if ms.Grid[ny][nx].IsMine {
							count++
						}
					}
				}
			}
			ms.Grid[y][x].Adjacent = count
		}
	}
}

func (ms *Minesweeper) Reveal(x, y int) bool {
	if x < 0 || x >= ms.Width || y < 0 || y >= ms.Height {
		return false
	}

	cell := &ms.Grid[y][x]

	if cell.State == CellRevealed || cell.State == CellFlagged {
		return false
	}

	if cell.IsMine {
		ms.GameOver = true
		return true
	}

	ms.revealCell(x, y)
	ms.checkWin()
	return true
}

func (ms *Minesweeper) revealCell(x, y int) {
	if x < 0 || x >= ms.Width || y < 0 || y >= ms.Height {
		return
	}

	cell := &ms.Grid[y][x]

	if cell.State == CellRevealed || cell.State == CellFlagged {
		return
	}

	cell.State = CellRevealed
	ms.Revealed++

	if cell.Adjacent == 0 {
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				if dx == 0 && dy == 0 {
					continue
				}
				ms.revealCell(x+dx, y+dy)
			}
		}
	}
}

func (ms *Minesweeper) ToggleFlag(x, y int) {
	if x < 0 || x >= ms.Width || y < 0 || y >= ms.Height {
		return
	}

	cell := &ms.Grid[y][x]

	if cell.State == CellRevealed {
		return
	}

	if cell.State == CellFlagged {
		cell.State = CellHidden
		ms.Flags--
	} else {
		cell.State = CellFlagged
		ms.Flags++
	}
}

func (ms *Minesweeper) checkWin() {
	totalCells := ms.Width * ms.Height
	if ms.Revealed == totalCells-ms.MineCount {
		ms.Won = true
		ms.GameOver = true
	}
}

func (ms *Minesweeper) GetElapsedTime() time.Duration {
	if ms.GameOver {
		return time.Since(ms.StartTime)
	}
	return time.Since(ms.StartTime)
}
