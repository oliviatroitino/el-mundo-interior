package handlers

import (
	"el-mundo-interior/internal/content"
	"log"
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

		allPosts, err := posts.GetBySection(worldSlug, sectionSlug)
		if err != nil {
			log.Printf("error cargando posts de sección %s/%s: %v", worldSlug, sectionSlug, err)
		}

		currentUserID, userName, loggedIn := sessions.GetUser(r)
		var myPosts, otherPosts []Post
		for _, p := range allPosts {
			if loggedIn && p.UserID == currentUserID {
				myPosts = append(myPosts, toViewPost(p))
			} else {
				otherPosts = append(otherPosts, toViewPost(p))
			}
		}

		data := SectionPageData{
			World:   world,
			Section: section,
			Nav: NavData{
				HomeHref:     "/",
				NavDropdowns: []NavDropdown{buildWorldDropdown(worldSlug), buildSectionDropdown(world.Sections, sectionSlug)},
				UserDropdown: func() *NavDropdown { ud := buildUserDropdown(userName); return &ud }(),
			},
			Questions: []string{
				"¿Qué emoción quieres transmitir con esta expresión?",
				"¿Qué momento de la realidad estás eligiendo capturar y por qué merece ser observado?",
				"¿Qué historia puede entenderse sin necesidad de palabras?",
				"¿Qué herramienta o técnica te permitiría expresar mejor la idea que tienes ahora?",
			},
			MyPosts:    myPosts,
			OtherPosts: otherPosts,
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

		r.ParseMultipartForm(10 << 20)
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
