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
		Questions: []string{
			"¿Qué emoción quieres transmitir con esta expresión?",
			"¿Qué momento de la realidad estás eligiendo capturar y por qué merece ser observado?",
			"¿Qué historia puede entenderse sin necesidad de palabras?",
			"¿Qué herramienta o técnica te permitiría expresar mejor la idea que tienes ahora?",
		},
		MyPosts: []Post{
			{User: "tu", Title: "Mi obra reciente", Text: "Descripción de mi obra.", Location: "Madrid", Date: "2025-11-01"},
		},
		OtherPosts: []Post{
			{User: "usuario1", Title: "Nueva obra", Text: "Esta es una breve descripción de la obra compartida.", Location: "Buenos Aires", Date: "2025-11-01"},
			{User: "usuario2", Title: "Inspiración nocturna", Text: "Reflexión sobre el proceso creativo.", Location: "Madrid", Date: "2025-11-03"},
		},
	}

	render(w, "templates/pages/section.tmpl", data)
}
