package main

import (
	"html/template"
	"log"
	"net/http"
)

// PageData holds the data we pass into the template.
type PageData struct {
	Title string
	Name  string
}

// PageData holds the data we pass into the template.
type Product struct {
	Name  string
	Price int
}

type ProductPage struct {
	Title    string
	Products []Product
}

var templates = template.Must(template.ParseGlob("./week04/lesson03/demo/templates/*.html"))

func productsHandler(w http.ResponseWriter, r *http.Request) {

	data := ProductPage{
		Title: "Danh sách sản phẩm",
		Products: []Product{
			{Name: "Trà sữa trân châu", Price: 30000},
			{Name: "Cà phê muối", Price: 25000},
			{Name: "Bánh mì Go", Price: 15000},
		},
	}
	if err := templates.ExecuteTemplate(w, "products.html", data); err != nil {
		log.Println(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{Title: "Trang chủ", Name: "bạn"}
	// Note: pass "." into {{template "content" .}} so content sees the same data.
	if err := templates.ExecuteTemplate(w, "layout", data); err != nil {
		log.Println(err)
	}
}

func main() {
	// Serve everything under ./static at URL path /static/
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)
	log.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
