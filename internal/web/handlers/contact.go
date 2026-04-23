package handlers

import (
	"net/http"
	"strings"
)

// Contact handles GET and POST /contacto.
func Contact() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			render(w, "templates/pages/contact.tmpl", ContactPageData{
				Nav: NavData{
					HomeHref: "/",
				},
				Success: r.URL.Query().Get("ok") == "1",
			})
			return
		}

		name := strings.TrimSpace(r.FormValue("name"))
		email := strings.TrimSpace(r.FormValue("email"))
		message := strings.TrimSpace(r.FormValue("message"))

		if name == "" || email == "" || message == "" {
			render(w, "templates/pages/contact.tmpl", ContactPageData{
				Nav: NavData{
					HomeHref: "/",
				},
				Error: "Todos los campos son obligatorios.",
			})
			return
		}

		http.Redirect(w, r, "/contacto?ok=1", http.StatusSeeOther)
	}
}
