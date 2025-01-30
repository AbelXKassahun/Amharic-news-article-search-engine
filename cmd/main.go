package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"search_engine/searchEngine"
	"search_engine/utils"
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
	mux.HandleFunc("/details/", detailsHandler)

	log.Print("Listening on port 8000")
	err := http.ListenAndServe(":8000", mux)

	log.Fatalf("couldn't listen to port 8000: %v \n", err)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}
/*
type Article struct {
	Article_ID    string      `json:"article_id"`
	Title         string      `json:"title"`
	PublishedTime time.Time   `json:"published_time"`
	ArticleImage  image_info  `json:"article_image"`
	ArticleMeta   ArticleMeta `json:"article_meta"`
	Content       []Content   `json:"content"`
}
*/
var Articles []utils.Article
func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("query")
	if query != "" {
		// Add your search logic here (filter items based on query)
		items = generateSampleData()
		Articles = searchEngine.SearchEngine(query)
		tmpl := template.Must(template.ParseFiles("templates/cards.html"))
		tmpl.Execute(w, Articles) // Return first 10 items
	}
}

func detailsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
    articleID := strings.TrimPrefix(r.URL.Path, "/details/")

    // Find item by ID
    var articleDetail utils.Article
    for _, article := range Articles {
        if  article.Article_ID == articleID {
            articleDetail = article
            break
        }
    }

    if articleDetail.Article_ID == "" {
        http.NotFound(w, r)
        return
    }

    tmpl := template.Must(template.ParseFiles("templates/details.html"))
    tmpl.Execute(w, articleDetail)
}

func generateSampleData() []Item {
	var items []Item
	for i := 1; i <= 20; i++ {
		items = append(items, Item{
			ID:     i,
			Title:  fmt.Sprintf("Item %d", i),
			Image:  fmt.Sprintf("https://picsum.photos/300/200?random=%d", i),
			Detail: fmt.Sprintf("Detailed description for item %d", i),
		})
	}
	return items
}