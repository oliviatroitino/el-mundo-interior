package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"el-mundo-interior/internal/content"
)

// requireJSON rechaza peticiones cuyo Content-Type no sea application/json.
// Esto previene CSRF via formularios HTML (que no pueden enviar ese tipo).
func requireJSON(w http.ResponseWriter, r *http.Request) bool {
	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		writeError(w, http.StatusUnsupportedMediaType, "Content-Type debe ser application/json")
		return false
	}
	return true
}

// apiPost es la representación JSON de un post para la API.
type apiPost struct {
	ID          int    `json:"id"`
	UserName    string `json:"user_name"`
	WorldSlug   string `json:"world_slug"`
	SectionSlug string `json:"section_slug"`
	Body        string `json:"body"`
	Location    string `json:"location,omitempty"`
	MediaPath   string `json:"media_path,omitempty"`
	Date        string `json:"date"`
	Mine        bool   `json:"mine"`
}

// writeJSON serializa v como JSON y lo envía con el código de estado dado.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("writeJSON: error serializando respuesta: %v", err)
	}
}

// writeError envía un JSON {"error": msg} con el código de estado dado.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// toAPIPost convierte un content.Post en apiPost, marcando mine si es del usuario actual.
func toAPIPost(p content.Post, currentUserID int) apiPost {
	return apiPost{
		ID:          p.ID,
		UserName:    p.UserName,
		WorldSlug:   p.WorldSlug,
		SectionSlug: p.SectionSlug,
		Body:        p.Body,
		Location:    p.Location,
		MediaPath:   p.MediaPath,
		Date:        p.CreatedAt.Format("2006-01-02"),
		Mine:        p.UserID == currentUserID,
	}
}

// ApiGetPosts maneja GET /api/posts?world={slug}
// Devuelve todos los posts del mundo indicado como JSON array.
func ApiGetPosts(posts content.PostRepository, sessions *SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		worldSlug := r.URL.Query().Get("world")
		if worldSlug == "" {
			writeError(w, http.StatusBadRequest, "parámetro world requerido")
			return
		}
		sectionSlug := r.URL.Query().Get("section")

		var allPosts []content.Post
		var err error
		if sectionSlug != "" {
			allPosts, err = posts.GetBySection(worldSlug, sectionSlug)
		} else {
			allPosts, err = posts.GetByWorld(worldSlug)
		}
		if err != nil {
			log.Printf("error obteniendo posts: %v", err)
			writeError(w, http.StatusInternalServerError, "error obteniendo posts")
			return
		}

		currentUserID, _, _ := sessions.GetUser(r)

		result := make([]apiPost, len(allPosts))
		for i, p := range allPosts {
			result[i] = toAPIPost(p, currentUserID)
		}

		writeJSON(w, http.StatusOK, result)
	}
}

// ApiCreatePost maneja POST /api/posts
// Lee un JSON con world_slug, body, section_slug y location, crea el post y devuelve 201.
func ApiCreatePost(posts content.PostRepository, sessions *SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !requireJSON(w, r) {
			return
		}
		userID, userName, ok := sessions.GetUser(r)
		if !ok {
			writeError(w, http.StatusUnauthorized, "sesión requerida")
			return
		}

		var input struct {
			WorldSlug   string `json:"world_slug"`
			SectionSlug string `json:"section_slug"`
			Body        string `json:"body"`
			Location    string `json:"location"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		if input.Body == "" || input.WorldSlug == "" {
			writeError(w, http.StatusBadRequest, "body y world_slug son requeridos")
			return
		}

		id, err := posts.Create(content.Post{
			UserID:      userID,
			WorldSlug:   input.WorldSlug,
			SectionSlug: input.SectionSlug,
			Title:       deriveTitle(input.Body),
			Body:        input.Body,
			Location:    input.Location,
		})
		if err != nil {
			log.Printf("error creando post: %v", err)
			writeError(w, http.StatusInternalServerError, "error creando post")
			return
		}

		created := apiPost{
			ID:        id,
			UserName:  userName,
			WorldSlug: input.WorldSlug,
			Body:      input.Body,
			Location:  input.Location,
			Date:      time.Now().Format("2006-01-02"),
			Mine:      true,
		}

		writeJSON(w, http.StatusCreated, created)
	}
}

// ApiUpdatePost maneja PATCH /api/posts/{id}
// Lee un JSON con body, actualiza el post y devuelve el post modificado.
func ApiUpdatePost(posts content.PostRepository, sessions *SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !requireJSON(w, r) {
			return
		}
		userID, ok := sessions.GetUserID(r)
		if !ok {
			writeError(w, http.StatusUnauthorized, "sesión requerida")
			return
		}

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			writeError(w, http.StatusBadRequest, "id inválido")
			return
		}

		var input struct {
			Body        string `json:"body"`
			Location    string `json:"location"`
			SectionSlug string `json:"section_slug"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		if input.Body == "" {
			writeError(w, http.StatusBadRequest, "body requerido")
			return
		}

		if err := posts.Update(id, userID, input.Body, input.Location, input.SectionSlug); err != nil {
			if errors.Is(err, content.ErrNotOwner) {
				writeError(w, http.StatusForbidden, "no se puede actualizar el post")
			} else {
				log.Printf("error actualizando post %d: %v", id, err)
				writeError(w, http.StatusInternalServerError, "error actualizando post")
			}
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"id":           id,
			"body":         input.Body,
			"location":     input.Location,
			"section_slug": input.SectionSlug,
		})
	}
}

// ApiDeletePost maneja DELETE /api/posts/{id}
// Elimina el post indicado y devuelve 204 sin cuerpo.
func ApiDeletePost(posts content.PostRepository, sessions *SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := sessions.GetUserID(r)
		if !ok {
			writeError(w, http.StatusUnauthorized, "sesión requerida")
			return
		}

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			writeError(w, http.StatusBadRequest, "id inválido")
			return
		}

		if err := posts.Delete(id, userID); err != nil {
			if errors.Is(err, content.ErrNotOwner) {
				writeError(w, http.StatusForbidden, "no se puede borrar el post")
			} else {
				log.Printf("error borrando post %d: %v", id, err)
				writeError(w, http.StatusInternalServerError, "error borrando post")
			}
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
