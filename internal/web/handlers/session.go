package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
)

const sessionCookie = "session"

// SessionStore guarda las sesiones activas en memoria.
// Es un mapa de token aleatorio → userID.
// sync.RWMutex permite lecturas concurrentes seguras.
type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]int
}

// NewSessionStore crea un store vacío listo para usar.
func NewSessionStore() *SessionStore {
	return &SessionStore{sessions: make(map[string]int)}
}

// Create genera un token aleatorio, lo asocia al userID y lo devuelve.
func (s *SessionStore) Create(userID int) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)

	s.mu.Lock()
	s.sessions[token] = userID
	s.mu.Unlock()

	return token, nil
}

// GetUserID devuelve el userID de la sesión activa, o false si no hay sesión.
func (s *SessionStore) GetUserID(r *http.Request) (int, bool) {
	cookie, err := r.Cookie(sessionCookie)
	if err != nil {
		return 0, false
	}

	s.mu.RLock()
	id, ok := s.sessions[cookie.Value]
	s.mu.RUnlock()

	return id, ok
}

// SetCookie envía la cookie de sesión al navegador.
func (s *SessionStore) SetCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    token,
		Path:     "/",
		HttpOnly: true, // no accesible desde JavaScript
	})
}

// Clear elimina la sesión del store y borra la cookie del navegador.
func (s *SessionStore) Clear(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(sessionCookie); err == nil {
		s.mu.Lock()
		delete(s.sessions, cookie.Value)
		s.mu.Unlock()
	}
	http.SetCookie(w, &http.Cookie{
		Name:   sessionCookie,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}
