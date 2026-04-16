package users

import (
	"database/sql"
	"fmt"
)

// sqliteUserRepository implementa UserRepository usando SQLite.
type sqliteUserRepository struct {
	db *sql.DB
}

// NewUserRepository crea un repositorio de usuarios con la conexión dada.
func NewUserRepository(db *sql.DB) UserRepository {
	return &sqliteUserRepository{db: db}
}

// Create inserta un nuevo usuario y devuelve su ID.
func (r *sqliteUserRepository) Create(name, email, passwordHash string) (int, error) {
	result, err := r.db.Exec(
		`INSERT INTO users (name, email, password_hash) VALUES (?, ?, ?)`,
		name, email, passwordHash,
	)
	if err != nil {
		return 0, fmt.Errorf("creando usuario: %w", err)
	}
	id, err := result.LastInsertId()
	return int(id), err
}

// GetByEmail devuelve el usuario con ese email, o error si no existe.
func (r *sqliteUserRepository) GetByEmail(email string) (User, error) {
	var u User
	err := r.db.QueryRow(
		`SELECT id, name, email, password_hash, created_at FROM users WHERE email = ?`,
		email,
	).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return User{}, fmt.Errorf("buscando usuario: %w", err)
	}
	return u, nil
}
