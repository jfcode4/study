package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// helper function to log and return errors
func writeError(w http.ResponseWriter, msg string, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Printf("error: %s: %v", msg, err)
	fmt.Fprintf(w, "error: %s: %v", msg, err)
}

func handleDecks(api Api) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decks, err := json.Marshal(api.db.Decks)
		if err != nil {
			writeError(w, "failed to marshal json", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", decks)
	}
}
func handleDeck(api Api) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deckId, err := strconv.Atoi(r.PathValue("deckId"))
		if err != nil {
			http.Error(w, "deck id must be a number", http.StatusBadRequest)
			return
		}
		deck, err := api.db.GetDeck(deckId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		deckJson, err := json.Marshal(deck)
		if err != nil {
			writeError(w, "failed to marshal json", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", deckJson)
	}
}

func handleStudy(api Api) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deckId, err := strconv.Atoi(r.PathValue("deckId"))
		if err != nil {
			http.Error(w, "deck id must be a number", http.StatusBadRequest)
			return
		}
		cards, err := api.Study(deckId)
		if err != nil {
			writeError(w, "failed to get study cards", err)
			return
		}
		cardsJson, err := json.Marshal(cards)
		if err != nil {
			writeError(w, "failed to marshal json", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", cardsJson)

	}
}

func handleRate(api Api) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cardId, err1 := strconv.Atoi(r.PathValue("cardId"))
		rating, err2 := strconv.Atoi(r.PathValue("rating"))
		if err1 != nil || err2 != nil {
			http.Error(w, "cardId and rating must be integers", http.StatusBadRequest)
			return
		}
		err := api.Rate(cardId, rating)
		if err != nil {
			writeError(w, "failed to rate card", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `true`)
	}
}

func templateHandlers(mux *http.ServeMux, api Api) {
	tmpl := template.Must(template.ParseGlob("web/templates/*.html"))
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		decks := api.db.Decks
		tmpl.ExecuteTemplate(w, "home.html", decks)
	})
	mux.HandleFunc("GET /study/{deckId}", func(w http.ResponseWriter, r *http.Request) {
		deckId, err := strconv.Atoi(r.PathValue("deckId"))
		if err != nil {
			http.Error(w, "deck id must be a number", http.StatusBadRequest)
			return
		}
		cards, err := api.Study(deckId)
		if err != nil {
			writeError(w, "failed to get study cards", err)
			return
		}
		if len(cards) == 0 {
			http.Error(w, "there are no cards to study", http.StatusNotFound)
			return
		}
		cardsJson, err := json.Marshal(cards)
		if err != nil {
			writeError(w, "failed to marshal json", err)
			return
		}
		tmpl.ExecuteTemplate(w, "study.html", template.JS(cardsJson))
	})
}
