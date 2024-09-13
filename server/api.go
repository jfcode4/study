package main

import "time"

type Api struct {
	db *Database
}

func (api *Api) Study(id int) ([]Card, error) {
	today := time.Now().Add(time.Hour * -4).Format("2006-01-02")
	if today != api.db.CurrentDeck.Date {
		err := api.db.SetDeckDay(id, today, api.db.CurrentDeck.Day+1)
		if err != nil {
			return nil, err
		}
	}
	// get due cards
	cards, err := api.db.GetDueCards(id, today, 20)
	if err != nil {
		return cards, err
	}
	if len(cards) < 20 {
		// get new cards to fill 20
		rest := 20 - len(cards)
		newCards, err := api.db.GetNewCards(id, rest)
		if err != nil {
			return cards, err
		}
		cards = append(cards, newCards...)
	}
	return cards, nil
}

func (api *Api) Rate(cardId int, rating int) error {
	// get card
	card, err := api.db.GetCard(cardId)
	if err != nil {
		return err
	}
	today := time.Now().Add(time.Hour * -4).Format("2006-01-02")
	due := card.Due
	if due == "" {
		due = today
	}
	interval := card.Interval

	if rating == -1 {
		if interval != 0 {
			interval /= 3
		}
	} else if rating == 1 {
		if interval == 0 {
			interval = 1
		} else {
			interval *= 3
		}
	}
	if interval != 0 {
		dueTime, err := time.Parse("2006-01-02", due)
		if err != nil {
			return err
		}
		dueTime = dueTime.Add(time.Hour * time.Duration(24*interval))
		due = dueTime.Format("2006-01-02")
	}

	// update card
	err = api.db.SetCard(cardId, due, interval)
	if err != nil {
		return err
	}
	return nil
}
