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
	}

	for _, item := range content.OrderedWorlds() {
		data.Worlds = append(data.Worlds, WorldNavItem{
			Slug:   item.Slug,
			Title:  item.Title,
			Icon:   item.Icon,
			Active: item.Slug == worldSlug,
		})
	}

	render(w, "templates/pages/section.tmpl", data)
}
