// Package handlers contiene la función central de renderizado de templates.
package handlers

import (
	"errors"
	"html/template"
	"log"
	"net/http"
)

// templateBase son los archivos que se parsean siempre en todas las páginas:
// el layout base, y los partials compartidos (nav, footer, post card).
var templateBase = []string{
	"templates/layouts/base.tmpl",
	"templates/partials/nav.tmpl",
	"templates/partials/footer.tmpl",
	"templates/partials/post_card.tmpl",
	"templates/partials/planet_card.tmpl",
	"templates/partials/review_card.tmpl",
	"templates/partials/plan_card.tmpl",
}

// render parsea y ejecuta un template con los datos proporcionados.
// page es la ruta al archivo .tmpl de la página concreta.
// data es el view model (definido en viewmodels.go) que se pasa al template.
//
// Todos los handlers deben usar esta función en lugar de llamar a
// template.ParseFiles directamente.
func render(w http.ResponseWriter, page string, data any) {
	// Juntamos el layout base + partials + la página concreta.
	// Creamos un slice nuevo para no mutar templateBase.
	files := make([]string, len(templateBase)+1)
	copy(files, templateBase)
	files[len(templateBase)] = page

	tpl, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, "error cargando template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Ejecutamos el template "base", que es el punto de entrada definido
	// en templates/layouts/base.tmpl con {{define "base"}}
	if err := tpl.ExecuteTemplate(w, "base", data); err != nil {
		if !errors.Is(err, http.ErrAbortHandler) {
			log.Printf("error ejecutando template %s: %v", page, err)
			http.Error(w, "error renderizando página", http.StatusInternalServerError)
		}
	}
}
