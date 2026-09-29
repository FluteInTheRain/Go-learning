package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"Post/auth"

	"github.com/joho/godotenv"
)

type ctxKey string

const userClaimsKey ctxKey = "claims"

// writeJSON is a helper to send any value as a JSON response.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func loginHandler(db *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
			return
		}

		user := db.FindUserByEmail(req.Email)
		if user == nil || !auth.CheckPassword(req.Password, user.PasswordHash) {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "sai email hoặc mật khẩu"})
			return
		}

		token, err := auth.GenerateToken(user.ID, user.Role)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "lỗi máy chủ"})
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{"token": token})
	}
}

// AuthMiddleware rejects requests without a valid Bearer token.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		// Expected format: "Bearer <token>"
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "thiếu token", http.StatusUnauthorized)
			return
		}

		claims, err := auth.ParseToken(parts[1])
		if err != nil {
			http.Error(w, "token không hợp lệ hoặc hết hạn", http.StatusUnauthorized)
			return
		}

		// Attach claims to the request context so handlers can read who called.
		ctx := context.WithValue(r.Context(), userClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func meHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(userClaimsKey).(*auth.Claims)
	if !ok {
		http.Error(w, "không xác định được người dùng", http.StatusUnauthorized)
		return
	}
	fmt.Fprintf(w, "Xin chào user #%d với vai trò %s", claims.UserID, claims.Role)
}

func main() {
	godotenv.Load()

	db := newStore()

	// Create a test user for demo
	hashedPwd, _ := auth.HashPassword("password123")
	db.CreateUser("user@example.com", hashedPwd, "user")

	mux := http.NewServeMux()

	// POST /login -> authenticate and get JWT token
	mux.HandleFunc("POST /login", loginHandler(db))

	// GET /me -> get current user info (requires valid JWT token)
	mux.Handle("GET /me", AuthMiddleware(http.HandlerFunc(meHandler)))

	// GET /posts -> list all posts
	mux.HandleFunc("GET /posts", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, db.List())
	})

	// POST /posts -> create a new post
	mux.HandleFunc("POST /posts", func(w http.ResponseWriter, r *http.Request) {
		var p Post
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON không hợp lệ"})
			return
		}
		if p.Title == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Thiếu tiêu đề"})
			return
		}
		created := db.Create(p)
		writeJSON(w, http.StatusCreated, created)
	})

	// GET /posts/{id} -> get one post
	mux.HandleFunc("GET /posts/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id không hợp lệ"})
			return
		}
		p, ok := db.Get(id)
		if !ok {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "Không tìm thấy bài viết"})
			return
		}
		writeJSON(w, http.StatusOK, p)
	})

	// PUT /posts/{id} -> update a post
	mux.HandleFunc("PUT /posts/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id không hợp lệ"})
			return
		}
		var p Post
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON không hợp lệ"})
			return
		}
		updated, ok := db.Update(id, p)
		if !ok {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "Không tìm thấy bài viết"})
			return
		}
		writeJSON(w, http.StatusOK, updated)
	})

	// DELETE /posts/{id} -> delete a post
	mux.HandleFunc("DELETE /posts/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id không hợp lệ"})
			return
		}
		if !db.Delete(id) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "Không tìm thấy bài viết"})
			return
		}
		w.WriteHeader(http.StatusNoContent) // 204: no body
	})

	log.Println("Server chạy tại http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
