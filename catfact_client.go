package main

import (
	"encoding/json"
	"sync"
)

type CatfactClient struct {
	URL         string
	TimeoutTime float32
}

func (c *CatfactClient) Init() {
	c.TimeoutTime = 2
	c.URL = "https://cat-fact.herokuapp.com/"
}

func (c *CatfactClient) fetchNews(limit int, ch chan fetchedNews, wg *sync.WaitGroup) {
	defer wg.Done()

	body, err := getRequest(c.URL + "facts/")
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
