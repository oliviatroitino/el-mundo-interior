// Package handlers contiene las funciones que responden a las peticiones HTTP.
package handlers

import (
	"el-mundo-interior/internal/content"
	"errors"
	"html/template"
	"net/http"
)

type worldNavItem struct {
	Slug   string
	Title  string
	Icon   string
	Active bool
}

type worldPageData struct {
	Slug        string
	Title       string
	Description string
	Icon        string
	Worlds      []worldNavItem
	Sections    []content.WorldSection
}

// WorldBySlug maneja las rutas dinámicas /mundos/{slug}.
// Lee el slug de la URL, busca el mundo correspondiente y renderiza el template.
func WorldBySlug(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	slug := r.PathValue("slug")
	def, ok := content.GetWorldBySlug(slug)
	if !ok {
		http.NotFound(w, r)
		return
	}

	data := worldPageData{
		Slug:        slug,
		Title:       def.Title,
		Description: def.Description,
		Icon:        def.Icon,
		Sections:    def.Sections,
	}

	for _, item := range content.OrderedWorlds() {
		data.Worlds = append(data.Worlds, worldNavItem{
			Slug:   item.Slug,
			Title:  item.Title,
			Icon:   item.Icon,
			Active: item.Slug == slug,
		})
	}

	tpl, err := template.ParseFiles("templates/pages/world.tmpl")
	if err != nil {
		http.Error(w, "error loading template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tpl.Execute(w, data); err != nil {
		if !errors.Is(err, http.ErrAbortHandler) {
			http.Error(w, "error rendering page", http.StatusInternalServerError)
		}
		return
	}
}
