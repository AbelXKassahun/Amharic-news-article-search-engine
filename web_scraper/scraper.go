package web_scraper

// what's a corpus
// in what data format should i store the scraped data
// what should it look like for better use, (structure)
// how should i index

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type Content struct {
	ContentType string `json:"content_type"`
	Text        string `json:"text"`
}

type Article struct {
	Title         string    `json:"title"`
	PublishedTime time.Time `json:"published_time"`
	ArticleImage image_info `json:"article_image"`
	ArticleMeta ArticleMeta `json:"article_meta"`
	Content       []Content `json:"content"`
}

type image_info struct {
	Image_url string `json:"image_url"`
	Image_element string `json:"image_element"`
}

type ArticleMeta struct {
	ArticleUrl string `json:"articleUrl"`
	Image_meta image_info `json:"image_meta"`
}

func ScrapeArticleUrls(pageUrl string, numOfPages int) {
	// tells the collector which domains its allowed to scrape
	c := colly.NewCollector(
		colly.AllowedDomains("www.bbc.com", "bbc.com"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*bbc.com*",
		RandomDelay: 2 * time.Second,
		Parallelism: 1,
	})
	
	visitedPage := make(map[string]bool)
	for i := 0; i < numOfPages; i++ {
		articleMeta := []ArticleMeta{}

		// fired everytime we make a request
		c.OnRequest(func(r *colly.Request) {
			if _, exists := visitedPage[r.URL.String()]; exists {
				return
			}
			visitedPage[r.URL.String()] = true
			fmt.Printf("Visiting page %d, at %s\n", i+1, r.URL)
		})

		// called everytime there is an error
		c.OnError(func(r *colly.Response, e error) {
			fmt.Printf("Error while scraping: %s\n", e.Error())
		})
		
		// telling the collector what to do (calback function) when it meets a critera (html element using a selector)
		c.OnHTML("main li", func(h *colly.HTMLElement){
			img := h.DOM.Find("div div.promo-image div div img")
			imgElement, _ := goquery.OuterHtml(img)
			img_url, _ := img.Attr("src")

			if href, exists := h.DOM.Find("div.promo-text h2 a").Attr("href"); exists {
				articleMeta = append(articleMeta, ArticleMeta{
					ArticleUrl: href,
					Image_meta: image_info{
						Image_url: img_url,
						Image_element: strings.ReplaceAll(imgElement, "<nil>", ""),
					},
				})
			}

			// fmt.Println(href)
		})

		c.OnScraped(func(r *colly.Response) {
			_, err := json.MarshalIndent(articleMeta, "", "  ")
			if err != nil {
				fmt.Println("Error encoding JSON:", err)
				return
			}
			// fmt.Println(string(jsonData))
			// for j:=0; j < len(articleMeta); j++ {
			// 	ScrapeArticle(c, articleMeta[j], fmt.Sprintf("page %d",i+1))
			// }
			ScrapeArticle(c, articleMeta[0], fmt.Sprintf("page %d",i+1))
		})
		c.Visit(fmt.Sprintf("%s%d", pageUrl, i+1))
	}
}

func ScrapeArticle(c *colly.Collector, article_meta ArticleMeta, saveFileName string) Article{
	// scrapeUrl := "https://www.bbc.com/amharic/articles/cdjd3wj0nyro"
	article := Article{}

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting article %s\n", r.URL)
	})
	// div[dir='ltr'], main figure
	c.OnHTML("main", func(h *colly.HTMLElement) {
		selection := h.DOM
		content := selection.Find("div[dir='ltr'] > *")
		if content.Length() != 0 {
			article.ArticleMeta = article_meta

			img := selection.Find("figure div img")
			if img.Length() != 0 {
				imgElement, _ := goquery.OuterHtml(img)
				if img_url, exists := img.Attr("src"); exists {
					article.ArticleImage = image_info {
						Image_url: img_url,
						Image_element: imgElement,
					}
				}
			}
			content.Each(func(index int, element *goquery.Selection) {
				elementName := goquery.NodeName(element)
				switch elementName {
				case "h1":
					article.Title = element.Text()
				case "time":
					if publishedDate, exists := element.Attr("datetime"); exists {
						parsedTime, err := time.Parse("2006-01-02", publishedDate)
						if err != nil {
							log.Println("couldnt parse the published date")
						}
						article.PublishedTime = parsedTime
					}
					article.Content = append(article.Content, Content{
						ContentType: elementName,
						Text:        element.Text(),
					})
				case "p", "h2":
					article.Content = append(article.Content, Content{
						ContentType: elementName,
						Text:        element.Text(),
					})
					// fmt.Printf("content type: %s, text: %s\n", firstChildName, firstChild.Text())
				default:
					fmt.Println("No match found")
				}
			})
		}
	})

	// called when scraping is done
	c.OnScraped(func(r *colly.Response) {
		jsonData, err := json.MarshalIndent(article, "", "  ")
		if err != nil {
			fmt.Println("Error encoding JSON:", err)
			return
		}
		fmt.Println(string(jsonData))
		saveToJSON(saveFileName, article)
	})

	// tell the scraper to visit a specific url
	c.Visit(article_meta.ArticleUrl)

	return article
}

func saveToJSON(filename string, data Article) {
	file, err := os.Create("../corpus/"+filename+".json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Article saved to %s.json\n", filename)
}


// enc := json.NewEncoder(os.Stdout)
// enc.SetIndent("", " ")
// enc.Encode(news)