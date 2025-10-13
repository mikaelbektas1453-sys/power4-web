package main

import (
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

var game *Game

func main() {
	// Initialiser une nouvelle partie
	game = NewGame()

	// Définir les routes
	http.HandleFunc("/", showBoard)
	http.HandleFunc("/play", playMove)
	http.HandleFunc("/reset", resetGame)

	// Démarrer le serveur
	port := ":8082" // Changé de 8081 à 8082
	log.Printf("Power4 server starting on http://localhost%s", port)

	// Ouvrir le navigateur automatiquement après 1 seconde
	go func() {
		time.Sleep(1 * time.Second)
		openBrowser("http://localhost" + port)
	}()

	// Lancer le serveur
	log.Fatal(http.ListenAndServe(port, nil))
}

func showBoard(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	err := tmpl.Execute(w, game)
	if err != nil {
		log.Printf("Erreur template: %v", err)
		http.Error(w, "Erreur interne", http.StatusInternalServerError)
	}
}

func playMove(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	col, err := strconv.Atoi(r.FormValue("col"))
	if err != nil || col < 0 || col >= 7 || game.Board[0][col] != "" {
		log.Printf("Colonne invalide ou pleine: %v", col)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Jouer le coup
	game.Play(col)

	// Rediriger vers la page principale
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func resetGame(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Réinitialiser la partie
		game.Reset()
		log.Println("Nouvelle partie commencée")
	}

	// Rediriger vers la page principale
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		log.Printf("Système non supporté")
		return
	}

	if err != nil {
		log.Printf("Erreur ouverture navigateur: %v", err)
	} else {
		log.Println("Navigateur ouvert automatiquement")
	}
}
