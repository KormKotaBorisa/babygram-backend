package handlers

import (
	"babygram-backend/internal/database"
	"encoding/json"
	"net/http"
)

func GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	var name, email string
	err := database.DB.QueryRow("SELECT name, email FROM users WHERE id = $1", userID).Scan(&name, &email)
	if err != nil {
		http.Error(w, "User not found", 404)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"name":  name,
		"email": email,
	})
}