package web_scraper

// what's a corpus
// in what data format should i store the scraped data
// what should it look like for better use, (structure)
// how should i index

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"net/url"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type Content struct {
	ContentType string `json:"content_type"`
	Text        string `json:"text"`
}

type Article struct {
	Article_ID    string      `json:"article_id"`
	Title         string      `json:"title"`
	PublishedTime time.Time   `json:"published_time"`
	ArticleImage  image_info  `json:"article_image"`
	ArticleMeta   ArticleMeta `json:"article_meta"`
	Content       []Content   `json:"content"`
}

type image_info struct {
	Image_url     string `json:"image_url"`
	Image_element string `json:"image_element"`
}

type ArticleMeta struct {
	ArticleUrl string     `json:"articleUrl"`
	Image_meta image_info `json:"image_meta"`
}

func ScrapeArticleUrls(pageUrl string, numOfPages int, dirName string) {
	// tells the collector which domains its allowed to scrape
	c := colly.NewCollector(
		colly.AllowedDomains("www.bbc.com", "bbc.com"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"),
		// colly.Async(true),
		// colly.CacheDir("./cache"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*bbc.com*",
		RandomDelay: 2 * time.Second,
		Parallelism: 2,
	})

	visitedPage := make(map[string]bool)
	articleMeta := []ArticleMeta{}

	// fired everytime we make a request
	c.OnRequest(func(r *colly.Request) {
		if visitedPage[r.URL.String()] {
			r.Abort()
			return
		}
		visitedPage[r.URL.String()] = true
		fmt.Printf("Visiting page: %s\n", r.URL)
	})

	// called everytime there is an error
	c.OnError(func(r *colly.Response, e error) {
		fmt.Printf("Error while scraping: %s\n", e.Error())
	})

	for i := 0; i < 1; i++ {
		// telling the collector what to do (calback function) when it meets a critera (html element using a selector)
		c.OnHTML("main li", func(h *colly.HTMLElement) {
			img := h.DOM.Find("div div.promo-image div div img")
			imgElement, _ := goquery.OuterHtml(img)
			img_url, _ := img.Attr("src")

			if href, exists := h.DOM.Find("div.promo-text h2 a").Attr("href"); exists {
				href = normalizeURL(href)
				if !isDuplicate(articleMeta, href) {
					articleMeta = append(articleMeta, ArticleMeta{
						ArticleUrl: href,
						Image_meta: image_info{
							Image_url:     img_url,
							Image_element: strings.ReplaceAll(imgElement, "<nil>", ""),
						},
					})
				}
			}
			// fmt.Println(href)
		})

		c.OnScraped(func(r *colly.Response) {
			// commented code is to check the scraped page urls
			// _, err := json.MarshalIndent(articleMeta, "", "  ")
			// if err != nil {
			// 	fmt.Println("Error encoding JSON:", err)
			// 	return
			// }
			// fmt.Println(string(jsonData))
			var wg sync.WaitGroup
			for j := 0; j < len(articleMeta); j++ { // len(articleMeta)
				wg.Add(1)
				articleCollector := c.Clone()
				go func(meta ArticleMeta, page_num, article_num int) {
					defer wg.Done()
					// ScrapeArticle(articleCollector, meta, numOfPages+1, j+1, fmt.Sprintf("page_%d", j+1), dirName)
					ScrapeArticle(articleCollector, meta, page_num+1, article_num+1, fmt.Sprintf("page_%d", page_num+1), dirName)
				}(articleMeta[j], i, j)
			}
			wg.Wait()

			// ScrapeArticle(c, articleMeta[0], fmt.Sprintf("page %d", i+1), dirName)
		})

		// page_url := normalizeURL(fmt.Sprintf("%s%d", pageUrl, i+1))
		// fmt.Printf("%s\n%s\n",pageUrl, page_url)
		c.Visit(fmt.Sprintf("%s%d", pageUrl, i+1))
	}
	// c.Wait()
}

func isDuplicate(metaList []ArticleMeta, url string) bool {
	for _, m := range metaList {
		if m.ArticleUrl == url {
			return true
		}
	}
	return false
}

func ScrapeArticle(c *colly.Collector, article_meta ArticleMeta, page_num int, article_num int, saveFileName, dirName string) Article {
	// scrapeUrl := "https://www.bbc.com/amharic/articles/cdjd3wj0nyro"

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting article %s\n", r.URL)
	})

	article := Article{}

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
					article.ArticleImage = image_info{
						Image_url:     img_url,
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
		article_id := createArticleID(article_meta.ArticleUrl, page_num, article_num)
		article.Article_ID = article_id

		// commendted code is to verify scraped articles
		// jsonData, err := json.MarshalIndent(article, "", "  ")
		// if err != nil {
		// 	fmt.Println("Error encoding JSON:", err)
		// 	return
		// }
		// fmt.Println(string(jsonData))
		saveToJSON(dirName, saveFileName, article)
	})

	// tell the scraper to visit a specific url
	c.Visit(normalizeURL(article_meta.ArticleUrl))

	return article
}

func createArticleID(article_url string, page_num int, article_num int) string {
	parts := strings.Split(article_url, "/")
	var slug string
	if len(parts) >= 6 {
		slug = parts[5] // Assuming the identifier is at index 5
	} else {
		slug = "unknown" // Fallback for invalid URLs
	}

	// Format the ID (filesystem-safe: replace "#" with "_")
	return fmt.Sprintf("pg%d_ar%d_%s", page_num, article_num, slug)
}

func saveToJSON(dirname string, filename string, data Article) {
	filePath := filepath.Join("../corpus/", dirname, filename+".json")

	// Create directory if it doesn't exist
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		if err := os.MkdirAll(dirname, os.ModePerm); err != nil {
			log.Fatal(err)
		}
		log.Println("Directory created: ", dirname)
	} else if err != nil {
		log.Fatal("Error checking directory: ", err)
	}

	var arr []Article
	if _, err := os.Stat(filePath); err == nil {
		file, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatal("Error reading file: ", err)
		}

		// Unmarshal existing data into a slice
		if err := json.Unmarshal(file, &arr); err != nil {
			log.Fatal("Error decoding JSON: ", err)
		}
	}

	// Check if article already exists
	exists := false
	for _, a := range arr {
		if a.ArticleMeta.ArticleUrl == data.ArticleMeta.ArticleUrl {
			exists = true
			break
		}
	}

	if !exists {
		arr = append(arr, data)
		sort.Slice(arr, func(i, j int) bool {
			return arr[i].PublishedTime.After(arr[j].PublishedTime)
		})
		updatedJSON, err := json.MarshalIndent(arr, "", "  ")
		if err != nil {
			log.Fatalf("Error encoding JSON: %v\n", err)
			return
		}
		if err := os.WriteFile(filePath, updatedJSON, 0644); err != nil {
			log.Fatalf("Error writing file: %v\n", err)
			return
		}
		log.Println("File created and data written:", filePath)
	}
}

func normalizeURL(u string) string {
	parsed, err := url.Parse(u)
	if err != nil {
		return u
	}
	parsed.Fragment = ""                               // Remove fragments (e.g., #section)
	parsed.RawQuery = ""                               // Remove query parameters
	parsed.Path = strings.TrimSuffix(parsed.Path, "/") // Remove trailing slash
	// parsed.Path = "/" + parsed.Path // Ensure paths start with a slash
	return parsed.String()
}
