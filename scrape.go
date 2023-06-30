package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// main() contains code adapted from example found in Colly's docs:
// http://go-colly.org/docs/examples/basic/

type Post struct {
	Title string `json:"title"`
	Like  int    `json:"like"`
}

func main() {
	// Instantiate default collector
	c := colly.NewCollector()

	posts := []Post{}

	c.OnHTML("div.crayons-story", func(e *colly.HTMLElement) {
		e.ForEach("span.aggregate_reactions_counter", func(_ int, el *colly.HTMLElement) {
			var likeStr string
			if strings.Contains(el.Text, "reactions") {
				likeStr = strings.Trim(el.Text, " reactions")

			} else {
				likeStr = strings.Trim(el.Text, " reaction")
			}

			likeStr = strings.TrimSpace(likeStr)
			likeCount, err := strconv.Atoi(likeStr)

			if err != nil {
				fmt.Println("Failed to convert like count:", err)
				return
			}

			p := Post{
				Like:  likeCount,
				Title: e.ChildText("h2.crayons-story__title"),
			}
			posts = append(posts, p)
		})

		sort.Slice(posts, func(i, j int) bool {
			return posts[i].Like > posts[j].Like
		})

		saveToJSON(posts)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on https://dev.to
	c.Visit("https://dev.to/")
}

func saveToJSON(posts []Post) {

	jsonData, err := json.MarshalIndent(posts, "", "    ")
	fmt.Println(posts)
	if err != nil {
		fmt.Println("Failed to marshal reactions to JSON:", err)
		return
	}

	err = ioutil.WriteFile("output.json", jsonData, 0644)

	if err != nil {
		fmt.Println("Failed to write posts JSON file:", err)
		return
	}

	fmt.Println("Posts saved to output.json")
}
