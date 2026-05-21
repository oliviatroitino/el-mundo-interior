package handlers

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
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
// Si secret no es vacío, firma el valor de la cookie con HMAC-SHA256.
type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]sessionData
	secret   []byte
}

// NewSessionStore crea un store vacío listo para usar.
// secret se usa para firmar la cookie; puede estar vacío (sin firma).
func NewSessionStore(secret string) *SessionStore {
	return &SessionStore{
		sessions: make(map[string]sessionData),
		secret:   []byte(secret),
	}
}

// sign devuelve "token.HMAC" cuando hay secret configurado, o el token sin modificar.
func (s *SessionStore) sign(token string) string {
	if len(s.secret) == 0 {
		return token
	}
	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(token))
	return token + "." + hex.EncodeToString(mac.Sum(nil))
}

// verify extrae el token del valor firmado y comprueba la firma.
// Devuelve el token puro y true si es válido.
func (s *SessionStore) verify(signed string) (string, bool) {
	if len(s.secret) == 0 {
		return signed, true
	}
	parts := strings.SplitN(signed, ".", 2)
	if len(parts) != 2 {
		return "", false
	}
	token, gotMAC := parts[0], parts[1]
	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(token))
	wantMAC := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(gotMAC), []byte(wantMAC)) {
		return "", false
	}
	return token, true
}

// Create genera un token aleatorio, lo asocia al usuario y lo devuelve firmado.
func (s *SessionStore) Create(userID int, userName string) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)

	s.mu.Lock()
	s.sessions[token] = sessionData{userID: userID, userName: userName}
	s.mu.Unlock()

	return s.sign(token), nil
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

	token, ok := s.verify(cookie.Value)
	if !ok {
		return 0, "", false
	}

	s.mu.RLock()
	data, ok := s.sessions[token]
	s.mu.RUnlock()

	return data.userID, data.userName, ok
}

// SetCookie envía la cookie de sesión al navegador.
func (s *SessionStore) SetCookie(w http.ResponseWriter, signed string) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    signed,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

// Clear elimina la sesión del store y borra la cookie del navegador.
func (s *SessionStore) Clear(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(sessionCookie); err == nil {
		if token, ok := s.verify(cookie.Value); ok {
			s.mu.Lock()
			delete(s.sessions, token)
			s.mu.Unlock()
		}
	}
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}
