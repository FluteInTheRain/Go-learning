package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Response is a generic envelope for all API responses.
type HealthResponse struct {
	Status string `json:"status"`
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role`
}

var users = []User{
	{ID: 1, Name: "Lan Anh", Role: "admin"},
	{ID: 2, Name: "Minh Quan", Role: "member"},
	{ID: 3, Name: "Thu Ha", Role: "member"},
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	resp := HealthResponse{
		Status: "ok",
	}
	writeJSON(w, http.StatusOK, resp)
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	roleQuery := r.URL.Query().Get("role")

	filteredUsers := make([]User, 0)

	for _, user := range users {
		if roleQuery == "" || user.Role == roleQuery {
			filteredUsers = append(filteredUsers, user)
		}
	}

	writeJSON(w, http.StatusOK, filteredUsers)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "JSON không hợp lệ", http.StatusBadRequest)
		return
	}
	users = append(users, newUser)
	writeJSON(w, http.StatusCreated, newUser)
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUsersHandler(w, r)
	case http.MethodPost:
		createUserHandler(w, r)
	default:
		http.Error(w, "Method không được hỗ trợ", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/api/health", healthHandler)
	http.HandleFunc("/api/users", usersHandler)
	log.Println("Server chạy tại http://localhost:8080/api/health")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
