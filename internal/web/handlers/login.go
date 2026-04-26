package handlers

import (
	"database/sql"
	"el-mundo-interior/internal/users"
	"errors"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Login maneja GET y POST /login.
func Login(userRepo users.UserRepository, sessions *SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			render(w, "templates/pages/login.tmpl", LoginPageData{
				Nav: NavData{HomeHref: "/"},
			})
			return
		}

		// POST: leer campos del formulario
		email := r.FormValue("email")
		password := r.FormValue("password")

		user, err := userRepo.GetByEmail(email)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				log.Printf("login: error inesperado en BD: %v", err)
			}
			log.Printf("login: intento fallido (usuario no encontrado)")
			render(w, "templates/pages/login.tmpl", LoginPageData{
				Nav:   NavData{HomeHref: "/"},
				Error: "Email o contraseña incorrectos.",
				Email: email,
			})
			return
		}

		// Comparar la contraseña con el hash guardado
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
			log.Printf("login: intento fallido para userID=%d (contraseña incorrecta)", user.ID)
			render(w, "templates/pages/login.tmpl", LoginPageData{
				Nav:   NavData{HomeHref: "/"},
				Error: "Email o contraseña incorrectos.",
				Email: email,
			})
			return
		}

		// Crear sesión y enviar cookie al navegador
		token, err := sessions.Create(user.ID, user.Name)
		if err != nil {
			log.Printf("login: error creando sesión para userID=%d: %v", user.ID, err)
			http.Error(w, "error creando sesión", http.StatusInternalServerError)
			return
		}
		log.Printf("login: sesión iniciada para userID=%d", user.ID)
		sessions.SetCookie(w, token)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// Logout maneja POST /logout: borra la sesión y redirige a la home.
func Logout(sessions *SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessions.Clear(w, r)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
