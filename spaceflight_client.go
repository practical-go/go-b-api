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
	c.URL = "https://api.spaceflightnewsapi.net/v3/"
}

func (c *SpaceflightClient) fetchNews(limit int) ([]News, error) {
	url := fmt.Sprintf("%sarticles?_limit=%d", c.URL, limit)
	body, err := getRequest(url)
	if err != nil {
		return nil, err
	}

	var spfnews []News
	err = json.Unmarshal(body, &spfnews)
	if err != nil {
		return nil, err
	}

	return spfnews, nil
}
