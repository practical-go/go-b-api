package main

import "encoding/json"

type CatfactClient struct {
	URL         string
	TimeoutTime float32
}

func (c *CatfactClient) Init() {
	c.TimeoutTime = 2
	c.URL = "https://cat-fact.herokuapp.com/facts/"
}

func (c *CatfactClient) fetchCatFacts(limit int) ([]News, error) {
	body, err := getRequest(c.URL)
	if err != nil {
		return nil, err
	}

	var catFacts []CatFact
	err = json.Unmarshal(body, &catFacts)
	if err != nil {
		return nil, err
	}

	var news []News
	for i := 0; i < limit && i < len(catFacts); i++ {
		news = append(news, News{
			Title:   "Cat fact ^_^",
			Summary: catFacts[i].Text,
		})
	}

	return news, nil
}
