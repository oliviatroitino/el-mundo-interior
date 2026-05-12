package content

import (
	"database/sql"
	"fmt"
)

// sqlitePostRepository implementa PostRepository usando SQLite.
type sqlitePostRepository struct {
	db *sql.DB
}

// NewPostRepository crea un repositorio de posts con la conexión dada.
func NewPostRepository(db *sql.DB) PostRepository {
	return &sqlitePostRepository{db: db}
}

// GetByWorld devuelve todos los posts de un mundo con el nombre del autor.
func (r *sqlitePostRepository) GetByWorld(worldSlug string) ([]Post, error) {
	rows, err := r.db.Query(`
		SELECT p.id, p.user_id, u.name, p.world_slug, p.section_slug,
		       p.title, p.body, p.location, p.media_path, p.created_at
		FROM posts p
		JOIN users u ON u.id = p.user_id
		WHERE p.world_slug = ?
		ORDER BY p.created_at DESC
	`, worldSlug)
	if err != nil {
		return nil, fmt.Errorf("consultando posts: %w", err)
	}
	defer rows.Close()

	return scanPosts(rows)
}

// GetBySection devuelve los posts de una sección concreta.
func (r *sqlitePostRepository) GetBySection(worldSlug, sectionSlug string) ([]Post, error) {
	rows, err := r.db.Query(`
		SELECT p.id, p.user_id, u.name, p.world_slug, p.section_slug,
		       p.title, p.body, p.location, p.media_path, p.created_at
		FROM posts p
		JOIN users u ON u.id = p.user_id
		WHERE p.world_slug = ? AND p.section_slug = ?
		ORDER BY p.created_at DESC
	`, worldSlug, sectionSlug)
	if err != nil {
		return nil, fmt.Errorf("consultando posts de sección: %w", err)
	}
	defer rows.Close()

	return scanPosts(rows)
}

// Create inserta un nuevo post y devuelve su ID asignado.
func (r *sqlitePostRepository) Create(post Post) (int, error) {
	result, err := r.db.Exec(`
		INSERT INTO posts (user_id, world_slug, section_slug, title, body, location, media_path)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, post.UserID, post.WorldSlug, post.SectionSlug, post.Title, post.Body, post.Location, post.MediaPath)
	if err != nil {
		return 0, fmt.Errorf("creando post: %w", err)
	}
	id, err := result.LastInsertId()
	return int(id), err
}

// Update modifica los campos editables de un post solo si pertenece al usuario indicado.
func (r *sqlitePostRepository) Update(id, userID int, body, location, sectionSlug string) error {
	result, err := r.db.Exec(`
		UPDATE posts SET body = ?, title = SUBSTR(?, 1, 60), location = ?, section_slug = ?
		WHERE id = ? AND user_id = ?
	`, body, body, location, sectionSlug, id, userID)
	if err != nil {
		return fmt.Errorf("actualizando post: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("post %d no encontrado o no pertenece al usuario", id)
	}
	return nil
}

// Delete elimina un post solo si pertenece al usuario indicado.
func (r *sqlitePostRepository) Delete(id, userID int) error {
	result, err := r.db.Exec(`DELETE FROM posts WHERE id = ? AND user_id = ?`, id, userID)
	if err != nil {
		return fmt.Errorf("borrando post: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("post %d no encontrado o no pertenece al usuario", id)
	}
	return nil
}

// scanPosts lee las filas del resultado y las convierte en []Post.
func scanPosts(rows *sql.Rows) ([]Post, error) {
	var posts []Post
	for rows.Next() {
		var p Post
		var location, mediaPath sql.NullString
		if err := rows.Scan(
			&p.ID, &p.UserID, &p.UserName, &p.WorldSlug, &p.SectionSlug,
			&p.Title, &p.Body, &location, &mediaPath, &p.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("leyendo post: %w", err)
		}
		p.Location = location.String
		p.MediaPath = mediaPath.String
		posts = append(posts, p)
	}
	return posts, rows.Err()
}
