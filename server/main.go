package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	// connect to database
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	db, err := OpenDB(home + "/.local/share/study.db")
	if err != nil {
		log.Fatal(err)
	}

	api := Api{&db}

	mux := http.NewServeMux()
	// API endpoints
	mux.HandleFunc("GET /api/decks/", handleDecks(api))
	mux.HandleFunc("GET /api/deck/{deckId}", handleDeck(api))
	mux.HandleFunc("GET /api/study/{deckId}", handleStudy(api))
	mux.HandleFunc("GET /api/rate/{cardId}/{rating}", handleRate(api))
	mux.Handle("GET /", http.FileServer(http.Dir("web/static")))
	// templates
	templateHandlers(mux, api)

	// start server
	println("Running on http://127.0.0.1:8000")
	log.Fatalln(http.ListenAndServe(":8000", mux))
	db.Close()
}
