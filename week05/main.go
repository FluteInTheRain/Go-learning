package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // Postgres driver
)

func main() {
	db, err := setupDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := &ProductRepo{db: db}
	h := &ProductHandler{repo: repo}

	r := mux.NewRouter()
	r.HandleFunc("/products", h.Create).Methods(http.MethodPost)
	r.HandleFunc("/products", h.List).Methods(http.MethodGet)
	r.HandleFunc("/products/count", h.Count).Methods(http.MethodGet)
	r.HandleFunc("/products/{id}", h.Get).Methods(http.MethodGet)
	r.HandleFunc("/products/{id}", h.Update).Methods(http.MethodPut)
	r.HandleFunc("/products/{id}", h.PartialUpdate).Methods(http.MethodPatch)
	r.HandleFunc("/products/{id}", h.Delete).Methods(http.MethodDelete)

	log.Println("Server đang chạy tại http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
