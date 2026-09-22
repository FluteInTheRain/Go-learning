package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
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

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "pong")
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Bay gio la %s", time.Now().Format("15:04:33"))
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	aStr := r.URL.Query().Get("a")
	bStr := r.URL.Query().Get("b")

	if aStr == "" || bStr == "" {
		http.Error(w, "a va b la bat buoc", http.StatusBadRequest)
		return
	}

	a, err := strconv.Atoi(aStr)
	if err != nil {
		http.Error(w, "a phai la so", http.StatusBadRequest)
		return
	}

	b, err := strconv.Atoi(bStr)
	if err != nil {
		http.Error(w, "b phai la so", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%d + %d = %d\n", a, b, a+b)
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", pingHandler)
	mux.HandleFunc("/time", timeHandler)
	mux.HandleFunc("/calc", calculateHandler)
	fmt.Println("Nghe tai http://localhost:3000")
	http.ListenAndServe(":3000", mux)
}
