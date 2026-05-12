package handlers

import (
	"el-mundo-interior/internal/content"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// WorldBySlug maneja GET /mundos/{slug}.
func WorldBySlug(posts content.PostRepository, sessions *SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")

		def, ok := content.GetWorldBySlug(slug)
		if !ok {
			http.NotFound(w, r)
			return
		}

		_, userName, _ := sessions.GetUser(r)

		data := WorldPageData{
			Slug:        slug,
			Title:       def.Title,
			Description: def.Description,
			Icon:        def.Icon,
			Sections:    def.Sections,
			Nav: NavData{
				HomeHref:     "/",
				NavDropdowns: []NavDropdown{buildWorldDropdown(slug), buildSectionDropdown(def.Sections, "")},
				UserDropdown: func() *NavDropdown { ud := buildUserDropdown(userName); return &ud }(),
			},
		}

		render(w, "templates/pages/world.tmpl", data)
	}
}

// CreatePost maneja POST /mundos/{slug}.
// Lee el formulario, verifica la sesión y guarda el post en la BD.
func CreatePost(posts content.PostRepository, sessions *SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")

		userID, ok := sessions.GetUserID(r)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		r.ParseMultipartForm(10 << 20) // 10 MB máx
		body := r.FormValue("body")
		sectionSlug := r.FormValue("section_slug")

		if body == "" {
			http.Redirect(w, r, "/mundos/"+slug, http.StatusSeeOther)
			return
		}

		mediaPath, err := saveUpload(r, "media")
		if err != nil {
			log.Printf("error guardando archivo: %v", err)
		}

		_, err = posts.Create(content.Post{
			UserID:      userID,
			WorldSlug:   slug,
			SectionSlug: sectionSlug,
			Title:       deriveTitle(body),
			Body:        body,
			Location:    r.FormValue("location"),
			MediaPath:   mediaPath,
		})
		if err != nil {
			log.Printf("error creando post: %v", err)
		}

		http.Redirect(w, r, "/mundos/"+slug, http.StatusSeeOther)
	}
}

// saveUpload guarda el archivo del campo fieldName en assets/uploads/ y
// devuelve la ruta pública relativa (p.ej. "/uploads/abc123.jpg").
// Devuelve "" si no se subió ningún archivo.
func saveUpload(r *http.Request, fieldName string) (string, error) {
	file, header, err := r.FormFile(fieldName)
	if err != nil {
		return "", nil // sin archivo adjunto
	}
	defer file.Close()

	if err := os.MkdirAll("static/assets/uploads", 0755); err != nil {
		return "", err
	}

	ext := filepath.Ext(header.Filename)
	name := uuid.NewString() + ext
	dst, err := os.Create(filepath.Join("static/assets/uploads", name))
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return "/uploads/" + name, nil
}

// deriveTitle extrae los primeros 60 caracteres del cuerpo como título.
func deriveTitle(body string) string {
	runes := []rune(body)
	if len(runes) <= 60 {
		return body
	}
	return string(runes[:60]) + "…"
}

