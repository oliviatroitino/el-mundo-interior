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
			HomeHref: "/",
			Dropdowns: []NavDropdown{
				buildWorldDropdown(slug),
				buildSectionDropdown(def.Sections, ""),
				buildUserDropdown(),
			},
		},
		MyPosts: []Post{
			{User: "tu", Title: "Nueva obra", Text: "Esta es una breve descripción de la obra compartida.", Location: "Valparaíso", Date: "2025-11-01"},
			{User: "tu", Title: "Inspiración en la naturaleza", Text: "Hoy capturé la esencia del bosque en mi lienzo, mezclando verdes vibrantes y sombras profundas para evocar tranquilidad.", Location: "Valparaíso", Date: "2025-11-02"},
		},
		OtherPosts: []Post{
			{User: "luna_pintora", Title: "Paisaje urbano nocturno", Text: "Exploré colores fríos y luz de farolas en la escena para transmitir serenidad en la noche.", Location: "Valparaíso", Date: "2025-11-01"},
			{User: "isa_creativa", Title: "Collage de texturas", Text: "Combiné papel, tinta y recortes para mostrar cómo las texturas despiertan nuevas emociones.", Location: "Ciudad de México", Date: "2025-11-03"},
		},
	}

	render(w, "templates/pages/world.tmpl", data)
}
