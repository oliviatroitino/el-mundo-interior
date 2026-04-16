// Package users contiene la lógica de dominio y acceso a datos de usuarios.
package users

import "time"

// User representa un usuario registrado en la aplicación.
type User struct {
	ID           int
	Name         string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

// UserRepository define las operaciones de acceso a datos para usuarios.
// Los handlers dependen de esta interfaz, no de la implementación concreta.
type UserRepository interface {
	// Create inserta un nuevo usuario y devuelve su ID asignado.
	Create(name, email, passwordHash string) (int, error)
	// GetByEmail busca un usuario por email; devuelve error si no existe.
	GetByEmail(email string) (User, error)
}
