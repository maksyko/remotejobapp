package remoteglobal

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/ievgen-ma/remotejobapp/feeds"
	"github.com/ievgen-ma/remotejobapp/models"
	"log"
	"net/http"
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
	config.host = "remoteglobal.com"
	config.url = fmt.Sprintf("https://%s/jm-ajax/get_listings", config.host)
	return &PublicFeed{
		config:   config,
		BaseFeed: feeds.NewBaseFeed(name),
	}
}

type Date struct {
	Html string `json:"html"`
}

func (feed *PublicFeed) Connect() {
	counter := 0
	response, err := http.Get(feed.config.url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()


	var data Date
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		log.Fatal(err)
	}


	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data.Html))
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("li").Each(func(i int, s *goquery.Selection) {
		if counter < feed.Limit() {
			href, exists := s.Find("a").Attr("href")
			if exists {
				u, err := url.Parse(href)
				if err != nil {
					log.Fatal(err)
				}

				job := feed.GetDocument(href)
				title := job.Find(".container .entry-title").Text()
				company := job.Find(".website").Text()
				apply, exists := job.Find(".application_details a").Attr("href")

				if exists {
					post := &models.Post{
						Path:     u.Path,
						Name:     feed.Name(),
						Host:     feed.config.host,
						Title:    strings.TrimSpace(title),
						Apply:    strings.TrimSpace(apply),
						Company:  strings.TrimSpace(company),
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
	})

}
