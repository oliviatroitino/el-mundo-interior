package handlers

import (
	"el-mundo-interior/internal/content"
	"log"
	"net/http"
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

		allPosts, err := posts.GetByWorld(slug)
		if err != nil {
			log.Printf("error cargando posts del mundo %s: %v", slug, err)
		}

		// Separar mis posts de los de otros según la sesión activa
		currentUserID, loggedIn := sessions.GetUserID(r)
		var myPosts, otherPosts []Post
		for _, p := range allPosts {
			if loggedIn && p.UserID == currentUserID {
				myPosts = append(myPosts, toViewPost(p))
			} else {
				otherPosts = append(otherPosts, toViewPost(p))
			}
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
			MyPosts:    myPosts,
			OtherPosts: otherPosts,
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

		body := r.FormValue("body")
		sectionSlug := r.FormValue("section_slug")

		if body == "" {
			http.Redirect(w, r, "/mundos/"+slug, http.StatusSeeOther)
			return
		}

		_, err := posts.Create(content.Post{
			UserID:      userID,
			WorldSlug:   slug,
			SectionSlug: sectionSlug,
			Title:       deriveTitle(body),
			Body:        body,
			Location:    r.FormValue("location"),
		})
		if err != nil {
			log.Printf("error creando post: %v", err)
		}

		http.Redirect(w, r, "/mundos/"+slug, http.StatusSeeOther)
	}
}

// deriveTitle extrae los primeros 60 caracteres del cuerpo como título.
func deriveTitle(body string) string {
	runes := []rune(body)
	if len(runes) <= 60 {
		return body
	}
	return string(runes[:60]) + "…"
}

// toViewPost convierte content.Post en el view model Post para el template.
func toViewPost(p content.Post) Post {
	return Post{
		User:     p.UserName,
		Title:    p.Title,
		Text:     p.Body,
		Location: p.Location,
		Date:     p.CreatedAt.Format("2006-01-02"),
	}
}

// toViewPosts convierte un slice de content.Post en []Post.
func toViewPosts(posts []content.Post) []Post {
	result := make([]Post, len(posts))
	for i, p := range posts {
		result[i] = toViewPost(p)
	}
	return result
}
