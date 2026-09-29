package main

import (
	"sync"
	"time"
)

// Post represents a single blog article.
type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

// User represents a user account.
type User struct {
	ID           int64  `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"` // never expose hash to client
	Role         string `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

// store is a thread-safe in-memory storage for posts and users.
type store struct {
	mu          sync.RWMutex
	posts       map[int]Post
	nextPostID  int
	users       map[int64]User
	nextUserID  int64
}

func newStore() *store {
	return &store{
		posts:      make(map[int]Post),
		nextPostID: 1,
		users:      make(map[int64]User),
		nextUserID: 1,
	}
}
