package main

import (
	"io"
	"net/http"
	"time"
)

type httpClient struct {
	client http.Client
}

func newHTTPClient(timeoutMs time.Duration) httpClient {
	return httpClient{
		client: http.Client{
			Timeout: timeoutMs * time.Millisecond,
		},
	}
}

func (c *httpClient) getRequest(url string) ([]byte, error) {
	response, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
