package handlers

import (
	"el-mundo-interior/internal/content"
	"log"
	"math/rand"
	"net/http"
)

// WorldSectionBySlug maneja GET /mundos/{slug}/{section}.
func WorldSectionBySlug(posts content.PostRepository, sessions *SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		worldSlug := r.PathValue("slug")
		sectionSlug := r.PathValue("section")

		world, section, ok := content.GetSectionBySlug(worldSlug, sectionSlug)
		if !ok {
			http.NotFound(w, r)
			return
		}

		_, userName, _ := sessions.GetUser(r)

		perm := rand.Perm(len(Questions))
		questions := make([]string, 4)
		for i := range questions {
			questions[i] = Questions[perm[i]]
		}

		data := SectionPageData{
			World:   world,
			Section: section,
			Nav: NavData{
				HomeHref:     "/",
				NavDropdowns: []NavDropdown{buildWorldDropdown(worldSlug), buildSectionDropdown(world.Sections, sectionSlug)},
				UserDropdown: func() *NavDropdown { ud := buildUserDropdown(userName); return &ud }(),
			},
			Questions: questions,
		}

		render(w, "templates/pages/section.tmpl", data)
	}
}

// CreateSectionPost maneja POST /mundos/{slug}/{section}.
func CreateSectionPost(posts content.PostRepository, sessions *SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		worldSlug := r.PathValue("slug")
		sectionSlug := r.PathValue("section")

		userID, ok := sessions.GetUserID(r)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, "error procesando formulario", http.StatusBadRequest)
			return
		}
		body := r.FormValue("body")
		if body == "" {
			http.Redirect(w, r, "/mundos/"+worldSlug+"/"+sectionSlug, http.StatusSeeOther)
			return
		}

		mediaPath, err := saveUpload(r, "media")
		if err != nil {
			log.Printf("error guardando archivo: %v", err)
		}

		_, err = posts.Create(content.Post{
			UserID:      userID,
			WorldSlug:   worldSlug,
			SectionSlug: sectionSlug,
			Title:       deriveTitle(body),
			Body:        body,
			Location:    r.FormValue("location"),
			MediaPath:   mediaPath,
		})
		if err != nil {
			log.Printf("error creando post en sección: %v", err)
		}

		http.Redirect(w, r, "/mundos/"+worldSlug+"/"+sectionSlug, http.StatusSeeOther)
	}
}
