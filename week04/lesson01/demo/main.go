package main

import (
	"fmt"
	"net/http"
)

func greetHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "ban"
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Xin chao, %s!\n", name)
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Gioi thieu ve chung toi")
	})
	fmt.Println("Nghe tai http://localhost:8080")
	http.ListenAndServe(":8080", mux)
	// http.HandleFunc("/greet", greetHandler)
	// fmt.Println("Server dang chay tai localhost:3000")
	// fmt.Println("Nghe tai http://localhost:8080/greet?name=Quynh")
	// http.ListenAndServe(":8080", nil)
}
