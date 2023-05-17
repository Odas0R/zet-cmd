package model

type Link struct {
	ZettelID  string `db:"zettel_id"`
	LinkID    string `db:"link_id"`
	CreatedAt Time   `db:"created_at"`
	UpdatedAt Time   `db:"updated_at"`
}
