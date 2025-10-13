package main

type Game struct {
	Board         [6][7]string
	CurrentPlayer string
	Winner        string
	Draw          bool
}

func NewGame() *Game { return &Game{CurrentPlayer: "R"} }

func (g *Game) Play(col int) {
	if g.Winner != "" || g.Draw {
		return
	}
	for row := 5; row >= 0; row-- {
		if g.Board[row][col] == "" {
			g.Board[row][col] = g.CurrentPlayer
			if g.checkWin(row, col) {
				g.Winner = g.CurrentPlayer
			}
			if g.checkDraw() {
				g.Draw = true
			} else {
				g.switchPlayer()
			}
			break
		}
	}
}

func (g *Game) switchPlayer() {
	if g.CurrentPlayer == "R" {
		g.CurrentPlayer = "Y"
	} else {
		g.CurrentPlayer = "R"
	}
}

func (g *Game) checkWin(row, col int) bool {
	player := g.Board[row][col]
	for _, d := range [][2]int{{0, 1}, {1, 0}, {1, 1}, {1, -1}} {
		if g.countConsecutive(row, col, d[0], d[1], player)+g.countConsecutive(row, col, -d[0], -d[1], player)+1 >= 4 {
			return true
		}
	}
	return false
}

func (g *Game) countConsecutive(row, col, dr, dc int, player string) int {
	count := 0
	for r, c := row+dr, col+dc; r >= 0 && r < 6 && c >= 0 && c < 7 && g.Board[r][c] == player; r, c = r+dr, c+dc {
		count++
	}
	return count
}

func (g *Game) checkDraw() bool {
	for _, cell := range g.Board[0] {
		if cell == "" {
			return false
		}
	}
	return true
}

func (g *Game) Reset() {
	g.Board = [6][7]string{}
	g.CurrentPlayer, g.Winner = "R", ""
	g.Draw = false
}
