package handlers

import (
	"el-mundo-interior/internal/users"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Register maneja GET y POST /registro.
func Register(userRepo users.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			render(w, "templates/pages/register.tmpl", RegisterPageData{
				Nav: NavData{HomeHref: "/"},
			})
			return
		}

		// POST: leer campos del formulario
		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Validación básica
		if name == "" || email == "" || password == "" {
			render(w, "templates/pages/register.tmpl", RegisterPageData{
				Nav:   NavData{HomeHref: "/"},
				Error: "Todos los campos son obligatorios.",
				Name:  name,
				Email: email,
			})
			return
		}

		// Hashear la contraseña antes de guardarla
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			render(w, "templates/pages/register.tmpl", RegisterPageData{
				Nav:   NavData{HomeHref: "/"},
				Error: "Error al procesar la contraseña.",
				Name:  name,
				Email: email,
			})
			return
		}

		_, err = userRepo.Create(name, email, string(hash))
		if err != nil {
			render(w, "templates/pages/register.tmpl", RegisterPageData{
				Nav:   NavData{HomeHref: "/"},
				Error: "El correo ya está registrado.",
				Name:  name,
				Email: email,
			})
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
