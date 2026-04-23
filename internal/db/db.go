// Package db gestiona la conexión con la base de datos SQLite.
package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite" // registra el driver "sqlite" en database/sql
)

// NewDB abre (o crea) el fichero SQLite en la ruta indicada
// y aplica el esquema inicial si las tablas no existen.
func NewDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("abriendo base de datos: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("comprobando conexión: %w", err)
	}

	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("aplicando esquema: %w", err)
	}

	return db, nil
}

// migrate crea las tablas si no existen.
// Se ejecuta siempre al arrancar; IF NOT EXISTS hace que sea seguro repetirlo.
func migrate(db *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id            INTEGER PRIMARY KEY AUTOINCREMENT,
			name          TEXT NOT NULL,
			email         TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			created_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS posts (
			id           INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id      INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			world_slug   TEXT NOT NULL,
			section_slug TEXT,
			title        TEXT NOT NULL,
			body         TEXT NOT NULL,
			location     TEXT,
			created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS contact_messages (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			name       TEXT NOT NULL,
			email      TEXT NOT NULL,
			message    TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}
	return nil
}
