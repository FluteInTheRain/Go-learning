package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

// In-memory fake data. In a real app this would come from a database.
var users = []User{
	{ID: 1, Name: "Lan Anh", Role: "admin"},
	{ID: 2, Name: "Minh Quan", Role: "member"},
	{ID: 3, Name: "Thu Ha", Role: "member"},
}

// Response is a generic envelope for all API responses.
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data"`
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := Response{
		Success: true,
		Data:    users,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Lỗi JSON", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/api/users", usersHandler)
	log.Println("Server chạy tại http://localhost:8080/api/users")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
