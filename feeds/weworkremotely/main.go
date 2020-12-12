package weworkremotely

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/ievgen-ma/remotejobapp/feeds"
	"github.com/ievgen-ma/remotejobapp/models"
	"log"
	"net/url"
	"strings"
)

type PublicFeedConfig struct {
	url  string
	host string
}

type PublicFeed struct {
	*feeds.BaseFeed
	config *PublicFeedConfig
}

func NewPublicFeed(name string) *PublicFeed {
	config := &PublicFeedConfig{}
	config.host = "https://weworkremotely.com"
	return &PublicFeed{
		config:   config,
		BaseFeed: feeds.NewBaseFeed(name),
	}
}

func (feed *PublicFeed) Connect() {
	counter := 0
	doc := feed.GetDocument(fmt.Sprintf("%s/categories/remote-programming-jobs#job-listings", feed.config.host))
	doc.Find(".jobs article ul").Children().Each(func(i int, s *goquery.Selection) {
		if counter < feed.Limit() {
			className, exists := s.Attr("class")
			if exists && className != "feature" {
				href, exists := s.Find("a").Attr("href")

				if exists && len(href) > 2 {
					job := feed.GetDocument(fmt.Sprintf("%s/%s", feed.config.host, href))
					title := job.Find(".listing-header-container h1").Text()
					company := job.Find(".company-card h2 a").Text()
					apply, exists := job.Find(".apply_tooltip a").Attr("href")

					if exists {
						u, err := url.Parse(href)
						if err != nil {
							log.Fatal(err)
						}

						post := &models.Post{
							Path:    u.Path,
							Name:    feed.Name(),
							Host:    feed.config.host,
							Title:   strings.TrimSpace(title),
							Apply:   strings.TrimSpace(apply),
							Company: strings.TrimSpace(company),
						}
						saved, err := feed.SavePost(post)
						if err != nil {
							log.Fatal(err)
						}
						if saved {
							log.Println(fmt.Sprintf("Post:%v saved successfully ", post))
							counter++
						}
					}
				}
			}
		}
	})
}
