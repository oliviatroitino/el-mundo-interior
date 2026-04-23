package contact

import (
	"database/sql"
	"fmt"
)

// Repository defines operations on contact messages.
type Repository interface {
	Save(name, email, message string) error
}

type sqliteRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &sqliteRepository{db: db}
}

func (r *sqliteRepository) Save(name, email, message string) error {
	_, err := r.db.Exec(
		`INSERT INTO contact_messages (name, email, message) VALUES (?, ?, ?)`,
		name, email, message,
	)
	if err != nil {
		return fmt.Errorf("guardando mensaje: %w", err)
	}
	return nil
}
