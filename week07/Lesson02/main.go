// main.go
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func main() {
	// Read the DB connection string from an environment variable.
	// Compose will inject this value for us.
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	// Retry connecting because the DB may still be starting up.
	var err error
	for i := 0; i < 10; i++ {
		pool, err = pgxpool.New(context.Background(), dsn)
		if err == nil {
			if pingErr := pool.Ping(context.Background()); pingErr == nil {
				break
			}
		}
		log.Printf("waiting for database... attempt %d", i+1)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}
	defer pool.Close()

	// Ensure the table exists.
	_, err = pool.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS todos (id SERIAL PRIMARY KEY, title TEXT NOT NULL)`)
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	http.HandleFunc("/todos", handleTodos)
	log.Println("server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleTodos(w http.ResponseWriter, r *http.Request) {
	rows, err := pool.Query(context.Background(), "SELECT id, title FROM todos ORDER BY id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Todo struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	}
	var todos []Todo
	for rows.Next() {
		var t Todo
		if scanErr := rows.Scan(&t.ID, &t.Title); scanErr != nil {
			http.Error(w, scanErr.Error(), http.StatusInternalServerError)
			return
		}
		todos = append(todos, t)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}
