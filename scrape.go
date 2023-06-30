package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// main() contains code adapted from example found in Colly's docs:
// http://go-colly.org/docs/examples/basic/

type post struct {
	title string
	like  int
}

func main() {
	// Instantiate default collector
	c := colly.NewCollector()

	posts := []post{}

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

			fmt.Println(e.ChildText("h2.crayons-story__title"))

			p := post{
				like:  likeCount,
				title: e.ChildText("h2.crayons-story__title"),
			}
			posts = append(posts, p)
		})

		sort.Slice(posts, func(i, j int) bool {
			return posts[i].like > posts[j].like
		})

		for _, r := range posts {
			fmt.Printf("Like count: %d\n", r.like)
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on https://dev.to
	c.Visit("https://dev.to/")
}
