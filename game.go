package main

// Structure principale du jeu Power4
type Game struct {
	Board         [6][7]string // Plateau 6x7 : "" = vide, "R" = rouge, "Y" = jaune
	CurrentPlayer string       // Joueur actuel : "R" ou "Y"
	Winner        string       // Gagnant : "", "R" ou "Y"
	Draw          bool         // Match nul : true/false
}

// Crée une nouvelle partie (Rouge commence toujours)
func NewGame() *Game {
	return &Game{CurrentPlayer: "R"}
}

// Joue un coup dans la colonne spécifiée
func (g *Game) Play(col int) {
	// Vérifie si le jeu est terminé
	if g.Winner != "" || g.Draw {
		return // Sort si quelqu'un a gagné ou match nul
	}

	// Trouve la case la plus basse disponible (gravité)
	for row := 5; row >= 0; row-- { // De bas en haut (5→0)
		if g.Board[row][col] == "" { // Case vide trouvée
			g.Board[row][col] = g.CurrentPlayer // Place le jeton

			// Vérifie s'il y a victoire
			if g.checkWin(row, col) {
				g.Winner = g.CurrentPlayer // Marque le gagnant
			}

			// Vérifie match nul ou change de joueur
			if g.checkDraw() {
				g.Draw = true // Plateau plein
			} else {
				g.switchPlayer() // Passe au joueur suivant
			}
			break // Sort de la boucle (jeton placé)
		}
	}
}

// Change de joueur (Rouge ↔ Jaune)
func (g *Game) switchPlayer() {
	if g.CurrentPlayer == "R" {
		g.CurrentPlayer = "Y" // Rouge → Jaune
	} else {
		g.CurrentPlayer = "R" // Jaune → Rouge
	}
}

// Vérifie si le dernier coup crée une victoire (4 alignés)
func (g *Game) checkWin(row, col int) bool {
	player := g.Board[row][col] // Récupère le joueur qui vient de jouer

	// Teste les 4 directions possibles
	for _, d := range [][2]int{
		{0, 1},  // Horizontal (←→)
		{1, 0},  // Vertical (↑↓)
		{1, 1},  // Diagonale ↘
		{1, -1}, // Diagonale ↙
	} {
		// Compte dans les 2 sens + le jeton actuel
		if g.countConsecutive(row, col, d[0], d[1], player)+ // Sens +
			g.countConsecutive(row, col, -d[0], -d[1], player)+ // Sens -
			1 >= 4 { // + jeton actuel = 4 alignés ?
			return true // Victoire !
		}
	}
	return false // Pas de victoire
}

// Compte les jetons consécutifs dans une direction
func (g *Game) countConsecutive(row, col, dr, dc int, player string) int {
	count := 0
	// Parcourt dans la direction donnée
	for r, c := row+dr, col+dc; // Position suivante
	r >= 0 && r < 6 &&          // Dans les limites lignes
		c >= 0 && c < 7 && // Dans les limites colonnes
				g.Board[r][c] == player; // Même joueur
	r, c = r+dr, c+dc { // Case suivante
		count++ // Compte ce jeton
	}
	return count
}

// Vérifie si c'est un match nul (plateau plein)
func (g *Game) checkDraw() bool {
	// Vérifie seulement la ligne du haut
	for _, cell := range g.Board[0] { // Ligne 0 (haut)
		if cell == "" {
			return false // Case vide = pas plein
		}
	}
	return true // Toute la ligne du haut est pleine = match nul
}

// Remet le jeu à zéro
func (g *Game) Reset() {
	g.Board = [6][7]string{}            // Plateau vide
	g.CurrentPlayer, g.Winner = "R", "" // Rouge recommence, pas de gagnant
	g.Draw = false                      // Pas de match nul
}
