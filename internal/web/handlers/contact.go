package handlers

import (
	"el-mundo-interior/internal/contact"
	"log"
	"net/http"
	"strings"
)

// Contact handles the footer contact form.
// GET /contacto → redirect to home (not a navigable page).
// POST /contacto → save message and return 200 OK (no redirect, JS handles the UI).
func Contact(repo contact.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		name := strings.TrimSpace(r.FormValue("name"))
		email := strings.TrimSpace(r.FormValue("email"))
		message := strings.TrimSpace(r.FormValue("message"))

		if name == "" || email == "" || message == "" {
			http.Error(w, "Todos los campos son obligatorios.", http.StatusBadRequest)
			return
		}

		if err := repo.Save(name, email, message); err != nil {
			log.Printf("contacto: error guardando mensaje: %v", err)
			http.Error(w, "No se pudo guardar el mensaje.", http.StatusInternalServerError)
			return
		}

		log.Printf("contacto: mensaje recibido correctamente")
		w.WriteHeader(http.StatusOK)
	}
}
