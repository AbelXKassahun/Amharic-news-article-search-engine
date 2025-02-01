package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"search_engine/searchEngine"
	"search_engine/utils"
	"strings"
)


var PORT = 8000

func main() {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/search", searchHandler)
	mux.HandleFunc("/details/", detailsHandler)

	log.Print("Listening on port ", PORT)
	err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), mux)

	log.Fatalf("couldn't listen to port %d: %v \n", PORT, err)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

var Articles []utils.Article
// var prevQuery string
func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("query")
	// if query != "" && query != prevQuery {
	// 	prevQuery = query
	// }
	temp_articles, found := searchEngine.SearchEngine(query)
	if !found {
		tmpl := template.Must(template.ParseFiles("templates/no_results.html"))
		err := tmpl.Execute(w, query) 
		if err != nil {
			http.Error(w, "Failed to render the no results template.", http.StatusInternalServerError)
		}
		return
	}
	Articles = temp_articles
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

// package main

// import (
// 	"search_engine/searchEngine"
// 	// "fmt"

// 	// "github.com/AbelXKassahun/Amharic-Stemmer/stemmer"
// )



// func main(){

// 	// fmt.Println(stemmer.Stem("ሆኗል"))
// 	searchEngine.TermWeighing()
// }