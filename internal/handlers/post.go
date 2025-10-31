package handlers

import (
	"babygram-backend/internal/database"
	"babygram-backend/internal/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	var req models.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	var postID int
	err := database.DB.QueryRow(
		"INSERT INTO posts(user_id, title, photo) VALUES($1, $2, $3) RETURNING id",
		userID, req.Title, req.Photo,
	).Scan(&postID)

	if err != nil {
		http.Error(w, "Failed to create post", 500)
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(map[string]int{"id": postID})
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(`
		SELECT p.id, p.user_id, p.title, p.photo, p.likes, p.created_at, u.name
		FROM posts p
		JOIN users u ON p.user_id = u.id
		ORDER BY p.created_at DESC LIMIT 20
	`)
	if err != nil {
		http.Error(w, "DB error", 500)
		return
	}
	defer rows.Close()

	posts := []models.Post{}
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Photo, &p.Likes, &p.CreatedAt, &p.UserName); err != nil {
			continue
		}
		posts = append(posts, p)
	}

	json.NewEncoder(w).Encode(posts)
}

func LikePost(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	postID, _ := strconv.Atoi(chi.URLParam(r, "id"))

	_, err := database.DB.Exec(
		"INSERT INTO likes(user_id, post_id) VALUES($1, $2) ON CONFLICT DO NOTHING",
		userID, postID,
	)
	if err != nil {
		http.Error(w, "Already liked", 400)
		return
	}

	database.DB.Exec("UPDATE posts SET likes = likes + 1 WHERE id = $1", postID)
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]string{"message": "Liked"})
}