package main

type Game struct {
	Board         [6][7]string // "R", "Y", or ""
	CurrentPlayer string       // "R" or "Y"
	Winner        string
	Draw          bool
}

func NewGame() *Game {
	return &Game{CurrentPlayer: "R"}
}

func (g *Game) Play(col int) {
	if g.Winner != "" || g.Draw {
		return
	}
	for row := 5; row >= 0; row-- {
		if g.Board[row][col] == "" {
			g.Board[row][col] = g.CurrentPlayer
			if g.checkWin(row, col) {
				g.Winner = g.CurrentPlayer
			} else if g.checkDraw() {
				g.Draw = true
			} else {
				g.switchPlayer()
			}
			break
		}
	}
}

