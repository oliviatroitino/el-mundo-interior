package handlers

import (
	"el-mundo-interior/internal/contact"
	"log"
	"net/http"
)

// Contact maneja el formulario de contacto del footer.
// GET /contacto → redirige a home (la página no es navegable directamente).
// POST /contacto → guarda el mensaje y redirige de vuelta.
func Contact(repo contact.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		name := r.FormValue("name")
		email := r.FormValue("email")
		message := r.FormValue("message")

		// En caso de error, volvemos a la página desde la que se envió el formulario.
		// Si no hay Referer (caso raro), mandamos a home.
		back := r.Referer()
		if back == "" {
			back = "/"
		}

		if name == "" || email == "" || message == "" {
			http.Redirect(w, r, back, http.StatusSeeOther)
			return
		}

		if err := repo.Save(name, email, message); err != nil {
			log.Printf("contacto: error guardando mensaje: %v", err)
			http.Redirect(w, r, back, http.StatusSeeOther)
			return
		}

		log.Printf("contacto: mensaje recibido correctamente")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
