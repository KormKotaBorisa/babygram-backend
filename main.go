package main

import (
	"babygram-backend/internal/database"
	"babygram-backend/internal/handlers"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Инициализация БД
	database.InitDB()
	defer database.CloseDB()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Healthcheck
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Публичные маршруты
	r.Post("/api/register", handlers.Register)
	r.Post("/api/login", handlers.Login)

	// Защищённые маршруты
	r.Group(func(r chi.Router) {
		r.Use(handlers.AuthMiddleware)
		r.Get("/api/profile", handlers.GetProfile)
		r.Post("/api/posts", handlers.CreatePost)
		r.Get("/api/posts", handlers.GetPosts)
		r.Post("/api/posts/{id}/like", handlers.LikePost)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Babygram API running on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}