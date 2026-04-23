// Package contact gestiona los mensajes enviados desde el formulario de contacto.
package contact

import (
	"database/sql"
	"fmt"
	"time"
)

// Message representa un mensaje de contacto guardado en la BD.
type Message struct {
	ID        int
	Name      string
	Email     string
	Body      string
	CreatedAt time.Time
}

// Repository define las operaciones sobre mensajes de contacto.
type Repository interface {
	Save(name, email, message string) error
}

// sqliteRepository implementa Repository usando SQLite.
type sqliteRepository struct {
	db *sql.DB
}

// NewRepository crea un repositorio de contacto con la conexión dada.
func NewRepository(db *sql.DB) Repository {
	return &sqliteRepository{db: db}
}

func (r *sqliteRepository) Save(name, email, message string) error {
	_, err := r.db.Exec(
		`INSERT INTO contact_messages (name, email, message) VALUES (?, ?, ?)`,
		name, email, message,
	)
	if err != nil {
		return fmt.Errorf("guardando mensaje de contacto: %w", err)
	}
	return nil
}
