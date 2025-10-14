package main

import (
	"html/template" // Pour traiter les templates HTML
	"log"           // Pour les messages de log (erreurs, infos)
	"net/http"      // Pour créer le serveur HTTP
	"os/exec"       // Pour exécuter des commandes système (ouvrir navigateur)
	"runtime"       // Pour détecter le système d'exploitation
	"strconv"       // Pour convertir string ↔ int
	"time"          // Pour les temporisations
)

// Variable globale qui stocke l'état du jeu
var game *Game

// FONCTION PRINCIPALE - Point d'entrée du programme
func main() {
	// Initialise une nouvelle partie
	game = NewGame() // Crée un nouveau jeu (Rouge commence)

	// Définit les routes HTTP (URLs que le serveur peut traiter)
	http.HandleFunc("/", showBoard)      // Page d'accueil → affiche le plateau
	http.HandleFunc("/play", playMove)   // Route pour jouer un coup
	http.HandleFunc("/reset", resetGame) // Route pour recommencer

	// Configuration du serveur
	port := ":8082" // Port d'écoute (localhost:8082)
	log.Printf("Power4 server starting on http://localhost%s", port)

	// Lance l'ouverture automatique du navigateur (en arrière-plan)
	go func() { // Goroutine (thread parallèle)
		time.Sleep(1 * time.Second)            // Attend 1 seconde
		openBrowser("http://localhost" + port) // Ouvre le navigateur
	}()

	// Démarre le serveur HTTP (BLOQUE ici jusqu'à arrêt)
	log.Fatal(http.ListenAndServe(port, nil)) // Si erreur → arrêt du programme
}

// AFFICHAGE DU PLATEAU - Gère la route "/"
func showBoard(w http.ResponseWriter, r *http.Request) {
	// Charge le template HTML depuis le fichier
	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	// Exécute le template avec les données du jeu
	if err := tmpl.Execute(w, game); err != nil { // w = réponse HTTP, game = données
		log.Printf("Erreur template: %v", err)                          // Log l'erreur
		http.Error(w, "Erreur interne", http.StatusInternalServerError) // Erreur 500
	}
}

// JOUER UN COUP - Gère la route "/play"
func playMove(w http.ResponseWriter, r *http.Request) {
	// Vérifie que c'est bien une requête POST (pas GET)
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther) // Redirige vers l'accueil
		return
	}

	// Récupère la colonne depuis le formulaire HTML
	colStr := r.FormValue("col")     // Récupère "col" du formulaire
	col, err := strconv.Atoi(colStr) // Convertit string → int

	// Vérifie que la colonne est valide
	if err != nil || col < 0 || col >= 7 { // Erreur conversion OU hors limites
		log.Printf("Colonne invalide: %s", colStr)
		http.Redirect(w, r, "/", http.StatusSeeOther) // Retour à l'accueil
		return
	}

	// Vérifie que la colonne n'est pas pleine
	if game.Board[0][col] != "" { // Si ligne du haut occupée
		log.Printf("Colonne %d est pleine", col)
		http.Redirect(w, r, "/", http.StatusSeeOther) // Retour à l'accueil
		return
	}

	// Joue le coup et redirige
	game.Play(col)                                // Appelle la logique du jeu
	http.Redirect(w, r, "/", http.StatusSeeOther) // Retour à l'accueil (actualise)
}

// RÉINITIALISER LE JEU - Gère la route "/reset"
func resetGame(w http.ResponseWriter, r *http.Request) {
	// Vérifie que c'est bien une requête POST
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Remet le jeu à zéro
	game.Reset()                             // Appelle la fonction Reset du jeu
	log.Println("Nouvelle partie commencée") // Log l'événement

	// Redirige vers l'accueil
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// OUVERTURE AUTOMATIQUE DU NAVIGATEUR
func openBrowser(url string) {
	var err error

	// Détecte le système d'exploitation et utilise la bonne commande
	switch runtime.GOOS { // GOOS = Go Operating System
	case "linux":
		err = exec.Command("xdg-open", url).Start() // Linux
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start() // Windows
	case "darwin":
		err = exec.Command("open", url).Start() // macOS
	default:
		log.Printf("Système non supporté pour l'ouverture automatique")
		return // Abandonne si système inconnu
	}

	// Gère les erreurs d'ouverture
	if err != nil {
		log.Printf("Erreur ouverture navigateur: %v", err)
	} else {
		log.Println("Navigateur ouvert automatiquement")
	}
}
