package handlers

import (
	"el-mundo-interior/internal/users"
	"log"
	"net/http"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

var (
	reEmail   = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	reUpper   = regexp.MustCompile(`[A-Z]`)
	reLower   = regexp.MustCompile(`[a-z]`)
	reDigit   = regexp.MustCompile(`[0-9]`)
	reSpecial = regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`)
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
		r.ParseForm()
		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Validación básica
		if _, ok := r.Form["privacy"]; !ok {
			render(w, "templates/pages/register.tmpl", RegisterPageData{
				Nav:   NavData{HomeHref: "/"},
				Error: "Debes aceptar la política de privacidad.",
				Name:  name,
				Email: email,
			})
			return
		}

		if name == "" || email == "" || password == "" {
			render(w, "templates/pages/register.tmpl", RegisterPageData{
				Nav:   NavData{HomeHref: "/"},
				Error: "Todos los campos son obligatorios.",
				Name:  name,
				Email: email,
			})
			return
		}

		if !reEmail.MatchString(email) {
			render(w, "templates/pages/register.tmpl", RegisterPageData{
				Nav:   NavData{HomeHref: "/"},
				Error: "El correo electrónico no es válido.",
				Name:  name,
				Email: email,
			})
			return
		}

		var pwdErr string
		switch {
		case len(password) < 8:
			pwdErr = "La contraseña debe tener al menos 8 caracteres."
		case !reUpper.MatchString(password):
			pwdErr = "La contraseña debe contener al menos una letra mayúscula."
		case !reLower.MatchString(password):
			pwdErr = "La contraseña debe contener al menos una letra minúscula."
		case !reDigit.MatchString(password):
			pwdErr = "La contraseña debe contener al menos un número."
		case !reSpecial.MatchString(password):
			pwdErr = "La contraseña debe contener al menos un carácter especial (!@#$%...)."
		}
		if pwdErr != "" {
			render(w, "templates/pages/register.tmpl", RegisterPageData{
				Nav:   NavData{HomeHref: "/"},
				Error: pwdErr,
				Name:  name,
				Email: email,
			})
			return
		}

		// Hashear la contraseña antes de guardarla
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("registro: error hasheando contraseña: %v", err)
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
			log.Printf("registro: error creando usuario en BD: %v", err)
			render(w, "templates/pages/register.tmpl", RegisterPageData{
				Nav:   NavData{HomeHref: "/"},
				Error: "El correo ya está registrado.",
				Name:  name,
				Email: email,
			})
			return
		}

		log.Printf("registro: nuevo usuario creado correctamente")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
