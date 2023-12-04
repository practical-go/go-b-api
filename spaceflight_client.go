package main

import (
	"encoding/json"
	"fmt"
	"sync"
)

type SpaceflightClient struct {
	URL        string
	HTTPClient httpClient
}

func (c *SpaceflightClient) Init() {
	c.URL = "https://api.spaceflightnewsapi.net/v3/"
	c.HTTPClient = newHTTPClient(750)
}

func (c *SpaceflightClient) fetchNews(limit int, ch chan fetchedNews, wg *sync.WaitGroup) {
	defer wg.Done()

	url := fmt.Sprintf("%sarticles?_limit=%d", c.URL, limit)
	body, err := c.HTTPClient.getRequest(url)
	if err != nil {
		ch <- fetchedNews{
			news: nil,
			err:  err,
		}
		return
	}

	var spfnews []News
	err = json.Unmarshal(body, &spfnews)
	if err != nil {
		ch <- fetchedNews{
			news: nil,
			err:  err,
		}
		return
	}

	ch <- fetchedNews{
		news: spfnews,
		err:  nil,
	}

	close(ch)
}
