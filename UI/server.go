package UI

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Item struct {
	ID     int
	Title  string
	Image  string
	Detail string
}

var items []Item // Populate this with your data

// var tmpl = template.Must(template.ParseFiles("../UI/index.html"))

// func outTest(w http.ResponseWriter, r * http.Request) {
// 	fmt.Fprintf(w,
// 	`<h1>OLA</h1>`,
// 	)
// }

func main() {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/search", searchHandler)
	mux.HandleFunc("/details", detailsHandler)

	log.Print("Listening on port 8080")
	err := http.ListenAndServe(":8080", mux)

	log.Fatalf("couldn't listen to port 8080: %v \n", err)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	// query := r.FormValue("query")
	// Add your search logic here (filter items based on query)

	tmpl := template.Must(template.ParseFiles("templates/cards.html"))
	tmpl.Execute(w, items[:10]) // Return first 10 items
}

func detailsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	id := strings.TrimPrefix(r.URL.Path, "/details/")

	// Find item by ID (replace with your lookup logic)
	var foundItem Item
	for _, item := range items {
		if strconv.Itoa(item.ID) == id {
			foundItem = item
			break
		}
	}

	tmpl := template.Must(template.ParseFiles("templates/details.html"))
	tmpl.Execute(w, foundItem)
}

func generateSampleData() []Item {
	var items []Item
	for i := 1; i <= 100; i++ {
		items = append(items, Item{
			ID:     i,
			Title:  fmt.Sprintf("Item %d", i),
			Image:  fmt.Sprintf("https://picsum.photos/300/200?random=%d", i),
			Detail: fmt.Sprintf("Detailed description for item %d", i),
		})
	}
	return items
}
