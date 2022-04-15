package model

import (
	"database/sql"
)

type Tag struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func (p *Tag) GetTag(db *sql.DB) error {
	return db.QueryRow("SELECT name FROM tags WHERE id=? AND status != 'deleted'",
		p.ID).Scan(&p.Name)
}

func (p *Tag) UpdateTag(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE tags SET name=? WHERE id=?",
			p.Name, p.ID)

	return err
}

func (p *Tag) DeleteTag(db *sql.DB) error {
	_, err := db.Exec("UPDATE tags SET status='deleted' WHERE id=?", p.ID)

	return err
}

func (p *Tag) CreateTag(db *sql.DB) error {
	res, err := db.Exec(
		"INSERT INTO tags(name) VALUES(?) ;",
		p.Name)

	if err != nil {
		return err
	}
	p.ID, err = res.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}

func ListTags(db *sql.DB, start, count int) ([]Tag, error) {
	rows, err := db.Query(
		"SELECT id, name FROM tags WHERE status != 'deleted' LIMIT ? OFFSET ?",
		count, start)

	if err != nil {
		// log.Println(err)
		return nil, err
	}

	defer rows.Close()

	tags := []Tag{}

	for rows.Next() {
		var t Tag
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}

	return tags, nil
}
