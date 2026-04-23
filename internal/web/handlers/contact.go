package handlers

import (
	"el-mundo-interior/internal/contact"
	"net/http"
)

// Contact maneja GET y POST /contacto.
func Contact(repo contact.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// ?sent=1 se añade tras un POST exitoso para mostrar el mensaje de confirmación
			render(w, "templates/pages/contact.tmpl", ContactPageData{
				Nav:     NavData{HomeHref: "/"},
				Success: r.URL.Query().Get("sent") == "1",
			})
			return
		}

		name := r.FormValue("name")
		email := r.FormValue("email")
		message := r.FormValue("message")

		if name == "" || email == "" || message == "" {
			render(w, "templates/pages/contact.tmpl", ContactPageData{
				Nav:   NavData{HomeHref: "/"},
				Error: "Todos los campos son obligatorios.",
			})
			return
		}

		if err := repo.Save(name, email, message); err != nil {
			render(w, "templates/pages/contact.tmpl", ContactPageData{
				Nav:   NavData{HomeHref: "/"},
				Error: "No se pudo enviar el mensaje. Inténtalo de nuevo.",
			})
			return
		}

		// POST-Redirect-GET: redirige para que refrescar no reenvíe el formulario
		http.Redirect(w, r, "/contacto?sent=1", http.StatusSeeOther)
	}
}
