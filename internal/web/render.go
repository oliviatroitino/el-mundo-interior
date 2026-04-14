package web

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

// RenderTemplate parsea y ejecuta un template, escribiendo el resultado en w
func RenderTemplate(w http.ResponseWriter, templateName string, data interface{}) error {
	// Construir la ruta a los templates
	// El . significa "desde donde se ejecuta el programa"
	baseDir := "templates"

	// Parsear el layout base y cualquier otro template
	// ExecuteTemplate buscará en todos los templates parseados
	tmpl, err := template.ParseGlob(filepath.Join(baseDir, "**", "*.tmpl"))
	if err != nil {
		return fmt.Errorf("error parsing templates: %w", err)
	}

	// Ejecutar el template específico solicitado
	// w es donde se escribe la salida
	// templateName es el nombre del template (por ejemplo "home" o "world")
	// data son los valores que pasamos al template
	err = tmpl.ExecuteTemplate(w, templateName, data)
	if err != nil {
		return fmt.Errorf("error executing template %s: %w", templateName, err)
	}

	return nil
}
