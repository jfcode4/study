package main

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db          *sqlx.DB
	Decks       []Deck
	CurrentDeck Deck
}

func OpenDB(filename string) (Database, error) {
	db := Database{}
	var err error
	db.db, err = sqlx.Open("sqlite3", filename)
	if err != nil {
		return db, fmt.Errorf("OpenDB: failed to open database '%s': %w", filename, err)
	}
	err = db.db.Select(&db.Decks, "SELECT * FROM Decks")
	if err != nil {
		return db, fmt.Errorf("OpenDB: failed to select table Decks: %w", err)
	}
	return db, nil
}

func (db *Database) GetDeck(deckId int) (Deck, error) {
	for _, deck := range db.Decks {
		if deck.Id == deckId {
			return deck, nil
		}
	}
	return Deck{}, fmt.Errorf("GetDeck: not found")
}

func (db *Database) SetDeckDay(deckId int, date string, day int) error {
	_, err := db.db.Exec("UPDATE Decks SET date=?, day=? WHERE id=?", date, day, deckId)
	if err != nil {
		return fmt.Errorf("SetDeckDay: %w", err)
	}
	db.CurrentDeck.Date = date
	db.CurrentDeck.Day = day
	return nil
}

func (db *Database) GetDueCards(deckId int, date string, limit int) ([]Card, error) {
	cards := []Card{}
	err := db.db.Select(&cards, "SELECT * FROM Cards WHERE deck_id=? AND due<=? AND due!='' ORDER BY due LIMIT ?", deckId, date, limit)
	if err != nil {
		return cards, fmt.Errorf("GetDueCards failed to select cards: %w", err)
	}
	return cards, nil
}
func (db *Database) GetNewCards(deckId int, limit int) ([]Card, error) {
	cards := []Card{}
	err := db.db.Select(&cards, "SELECT * FROM Cards WHERE deck_id=? AND due='' LIMIT ?", deckId, limit)
	if err != nil {
		return cards, fmt.Errorf("GetNewCards failed to select cards: %w", err)
	}
	return cards, nil
}
func (db *Database) GetCard(cardId int) (Card, error) {
	card := Card{}
	err := db.db.Get(&card, "SELECT * FROM Cards WHERE id=?", cardId)
	if err == sql.ErrNoRows {
		return card, fmt.Errorf("no card found with id=%d", cardId)
	}
	if err != nil {
		return card, fmt.Errorf("netCard: %w", err)
	}
	return card, nil
}
func (db *Database) SetCard(cardId int, due string, interval int) error {
	_, err := db.db.Exec("UPDATE Cards SET due=?, interval=? WHERE id=?", due, interval, cardId)
	return err
}

func (db *Database) Close() {
	db.db.Close()
}
