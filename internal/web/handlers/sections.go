// Package handlers contiene las funciones que responden a las peticiones HTTP.

package handlers

import (
	"el-mundo-interior/internal/content"
	"errors"
	"html/template"
	"net/http"
)


type sectionPageData struct {
	World      content.World
	Section    content.WorldSection
	SubSections []content.SubSection
	Worlds     []worldNavItem
}

// WorldSectionBySlug maneja rutas /mundos/{worldSlug}/{sectionSlug}
func WorldSectionBySlug(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	worldSlug := r.PathValue("slug")
	sectionSlug := r.PathValue("section")

	world, section, ok := content.GetSectionBySlug(worldSlug, sectionSlug)
	if !ok {
		http.NotFound(w, r)
		return
	}

	data := sectionPageData{
		World:      world,
		Section:    section,
		SubSections: section.SubSections,
	}

	for _, w := range content.OrderedWorlds() {
		data.Worlds = append(data.Worlds, worldNavItem{
			Slug:   w.Slug,
			Title:  w.Title,
			Icon:   w.Icon,
			Active: w.Slug == worldSlug,
		})
	}

	tpl, err := template.ParseFiles("templates/pages/section.tmpl")
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



