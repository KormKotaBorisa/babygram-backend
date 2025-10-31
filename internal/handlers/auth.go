package handlers

import (
	"babygram-backend/internal/database"
	"babygram-backend/internal/models"
	"babygram-backend/utils"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 10)

	var id int
	err := database.DB.QueryRow(
		"INSERT INTO users(email, password, name) VALUES($1, $2, $3) RETURNING id",
		req.Email, string(hashed), req.Name,
	).Scan(&id)

	if err != nil {
		http.Error(w, "Email already exists", 400)
		return
	}

	token, _ := utils.GenerateToken(id)
	user := models.User{ID: id, Email: req.Email, Name: req.Name}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.LoginResponse{Token: token, User: user})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	var user models.User
	err := database.DB.QueryRow(
		"SELECT id, email, name, password FROM users WHERE email = $1",
		req.Email,
	).Scan(&user.ID, &user.Email, &user.Name, &user.Password)

	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		http.Error(w, "Invalid credentials", 401)
		return
	}

	token, _ := utils.GenerateToken(user.ID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.LoginResponse{Token: token, User: user})
}