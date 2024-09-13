package main

type Card struct {
	Id       int    `json:"id"`
	DeckId   int    `json:"deck" db:"deck_id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Due      string `json:"due"`
	Interval int    `json:"interval"`
	Tags     string `json:"tags"`
}

type Deck struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Date  string `json:"date"`
	Day   int    `json:"day"`
	Owner int    `json:"owner"`
}
