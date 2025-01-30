package main

import (
	// "fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"search_engine/searchEngine"
	"search_engine/utils"
)

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
// var prevQuery string
func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("query")
	// if query != "" && query != prevQuery {
	// 	prevQuery = query
	// }
	Articles = searchEngine.SearchEngine(query)
	tmpl := template.Must(template.ParseFiles("templates/cards.html"))
	tmpl.Execute(w, Articles) 
}

var articleTmpl = template.Must(template.New("article").Funcs(template.FuncMap{
	"safeHTML": func(s string) template.HTML { return template.HTML(s) }, // Register function
}).Parse(`
<div class="details-container">

    <h1>{{.Title}}</h1>
	<img src="{{ .ArticleImage.Image_url }}" alt="{{ .Title }}">
    <time>{{.PublishedTime.Format "January 2, 2006"}}</time>

    <div class="details">
        {{range .Content}}
            {{printf "<%s>%s</%s>" .ContentType .Text .ContentType | safeHTML}}
        {{end}}
    </div>
</div>
`))

func detailsHandler(w http.ResponseWriter, r *http.Request) {
    articleID := strings.TrimPrefix(r.URL.Path, "/details/")

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

    articleTmpl.Execute(w, articleDetail)
}