package main

import (
	"encoding/json"
	"fmt"
)

type SpaceflightClient struct {
	URL         string
	TimeoutTime int
}

func (c *SpaceflightClient) Init() {
	c.TimeoutTime = 3
	c.URL = "https://api.spaceflightnewsapi.net/v3/articles?_limit=%d"
}

func (c *SpaceflightClient) fetchSpaceflightNews(limit int) ([]News, error) {
	url := fmt.Sprintf(c.URL, limit)
	body, err := getRequest(url)
	if err != nil {
		return nil, err
	}

	var spfnews []SpaceflightNews
	err = json.Unmarshal(body, &spfnews)
	if err != nil {
		return nil, err
	}

	var news []News
	for _, spfnew := range spfnews {
		news = append(news, News{
			Title:   spfnew.Title,
			Summary: spfnew.Summary,
		})
	}

	return news, nil
}
