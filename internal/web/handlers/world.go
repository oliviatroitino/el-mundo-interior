// Package handlers contiene las funciones que responden a las peticiones HTTP.
package handlers

import (
	"el-mundo-interior/internal/content"
	"net/http"
)

// WorldBySlug maneja las rutas dinámicas GET /mundos/{slug}.
// Lee el slug de la URL, busca el mundo correspondiente y renderiza el template.
func WorldBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	def, ok := content.GetWorldBySlug(slug)
	if !ok {
		http.NotFound(w, r)
		return
	}

	data := WorldPageData{
		Slug:        slug,
		Title:       def.Title,
		Description: def.Description,
		Icon:        def.Icon,
		Sections:    def.Sections,
		Nav: NavData{
			Dropdowns: []NavDropdown{
				buildWorldDropdown(slug),
				buildSectionDropdown(def.Sections, ""),
			},
			Links: authLinks(),
		},
	}

	render(w, "templates/pages/world.tmpl", data)
}
