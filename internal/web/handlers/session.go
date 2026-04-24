package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
)

const sessionCookie = "session"

type sessionData struct {
	userID   int
	userName string
}

// SessionStore guarda las sesiones activas en memoria.
// Es un mapa de token aleatorio → sessionData.
// sync.RWMutex permite lecturas concurrentes seguras.
type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]sessionData
}

// NewSessionStore crea un store vacío listo para usar.
func NewSessionStore() *SessionStore {
	return &SessionStore{sessions: make(map[string]sessionData)}
}

// Create genera un token aleatorio, lo asocia al usuario y lo devuelve.
func (s *SessionStore) Create(userID int, userName string) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)

	s.mu.Lock()
	s.sessions[token] = sessionData{userID: userID, userName: userName}
	s.mu.Unlock()

	return token, nil
}

// GetUserID devuelve el userID de la sesión activa, o false si no hay sesión.
func (s *SessionStore) GetUserID(r *http.Request) (int, bool) {
	id, _, ok := s.GetUser(r)
	return id, ok
}

// GetUser devuelve el userID y userName de la sesión activa.
func (s *SessionStore) GetUser(r *http.Request) (int, string, bool) {
	cookie, err := r.Cookie(sessionCookie)
	if err != nil {
		return 0, "", false
	}

	s.mu.RLock()
	data, ok := s.sessions[cookie.Value]
	s.mu.RUnlock()

	return data.userID, data.userName, ok
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
