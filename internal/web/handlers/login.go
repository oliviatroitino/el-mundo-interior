package handlers

import (
	"database/sql"
	"el-mundo-interior/internal/users"
	"errors"
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
			// sql.ErrNoRows significa que el email no existe; cualquier error muestra el mismo mensaje
			if !errors.Is(err, sql.ErrNoRows) {
				// error inesperado de BD
			}
			render(w, "templates/pages/login.tmpl", LoginPageData{
				Nav:   NavData{HomeHref: "/"},
				Error: "Email o contraseña incorrectos.",
				Email: email,
			})
			return
		}

		// Comparar la contraseña con el hash guardado
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
			render(w, "templates/pages/login.tmpl", LoginPageData{
				Nav:   NavData{HomeHref: "/"},
				Error: "Email o contraseña incorrectos.",
				Email: email,
			})
			return
		}

		// Crear sesión y enviar cookie al navegador
		token, err := sessions.Create(user.ID)
		if err != nil {
			http.Error(w, "error creando sesión", http.StatusInternalServerError)
			return
		}
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
