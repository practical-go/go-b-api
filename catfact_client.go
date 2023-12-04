package main

import (
	"encoding/json"
	"sync"
)

type CatfactClient struct {
	URL        string
	HTTPClient httpClient
}

func (c *CatfactClient) Init() {
	c.URL = "https://cat-fact.herokuapp.com/"
	c.HTTPClient = newHTTPClient(750)
}

func (c *CatfactClient) fetchNews(limit int, ch chan fetchedNews, wg *sync.WaitGroup) {
	defer wg.Done()

	body, err := c.HTTPClient.getRequest(c.URL + "facts/")
	if err != nil {
		ch <- fetchedNews{
			news: nil,
			err:  err,
		}
		return
	}

	var catFacts []CatFact
	err = json.Unmarshal(body, &catFacts)
	if err != nil {
		ch <- fetchedNews{
			news: nil,
			err:  err,
		}
		return
	}

	var news []News
	for i := 0; i < limit && i < len(catFacts); i++ {
		news = append(news, News{
			Title:   "Cat fact ^_^",
			Summary: catFacts[i].Text,
		})
	}
	ch <- fetchedNews{
		news: news,
		err:  nil,
	}
	close(ch)
}
