package models

type Post struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Title     string `json:"title"`
	Photo     string `json:"photo"`
	Likes     int    `json:"likes"`
	CreatedAt string `json:"created_at"`
	UserName  string `json:"user_name"`
}

type CreatePostRequest struct {
	Title string `json:"title"`
	Photo string `json:"photo"`
}