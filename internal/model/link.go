package model

type Link struct {
	From      string `db:"zettel_id"`
	To        string `db:"link_id"`
	CreatedAt Time   `db:"created_at"`
	UpdatedAt Time   `db:"updated_at"`
}
