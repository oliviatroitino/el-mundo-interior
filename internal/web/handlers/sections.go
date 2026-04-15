// Package handlers contiene las funciones que responden a las peticiones HTTP.
package handlers

import (
	"el-mundo-interior/internal/content"
	"net/http"
)

// WorldSectionBySlug maneja las rutas GET /mundos/{slug}/{section}.
// Lee ambos slugs de la URL, busca la sección y renderiza el template.
func WorldSectionBySlug(w http.ResponseWriter, r *http.Request) {
	worldSlug := r.PathValue("slug")
	sectionSlug := r.PathValue("section")

	world, section, ok := content.GetSectionBySlug(worldSlug, sectionSlug)
	if !ok {
		http.NotFound(w, r)
		return
	}

	data := SectionPageData{
		World:   world,
		Section: section,
		Nav: NavData{
			HomeHref: "/",
			Dropdowns: []NavDropdown{
				buildWorldDropdown(worldSlug),
				buildSectionDropdown(world.Sections, sectionSlug),
				buildUserDropdown(),
			},
		},
	}

	render(w, "templates/pages/section.tmpl", data)
}
