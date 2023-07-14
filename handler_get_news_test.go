package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type catClientMock struct {
	fetchNewsError error
	news           []News
}

func (c *catClientMock) fetchNews(limit int) ([]News, error) {
	return c.news, c.fetchNewsError
}

type spfClientMock struct {
	fetchNewsError error
	news           []News
}

func (c *spfClientMock) fetchNews(limit int) ([]News, error) {
	return c.news, c.fetchNewsError
}

func TestNewsHandler(t *testing.T) {
	tests := []struct {
		name           string
		catClient      newsFetcher
		spfClient      newsFetcher
		tag            string
		limit          int
		expectedStatus int
		expectedNews   []News
	}{
		{
			"Success",
			&catClientMock{
				news:           nil,
				fetchNewsError: nil,
			},
			&spfClientMock{
				news:           nil,
				fetchNewsError: nil,
			},
			"",
			10,
			http.StatusOK,
			nil,
		},
		{
			"Failure",
			&catClientMock{
				news:           nil,
				fetchNewsError: errors.New("Erorr "),
			},
			&spfClientMock{
				news:           nil,
				fetchNewsError: errors.New("Another error"),
			},
			"",
			10,
			http.StatusInternalServerError,
			nil,
		},
		{
			"Invalid Tag",
			&catClientMock{
				news: []News{
					{
						Title:   "Cat fact 1",
						Summary: "Cats are cool",
					},
				},
				fetchNewsError: nil,
			},
			&spfClientMock{
				news:           nil,
				fetchNewsError: nil,
			},
			"EOILNQWEQWEQWOEQWE",
			10,
			http.StatusOK,
			[]News{
				{Title: "Cat fact 1", Summary: "Cats are cool"},
			},
		},
		{
			"Invalid Limit",
			&catClientMock{
				news: []News{
					{
						Title:   "Cat fact 1",
						Summary: "Cats are cool",
					},
				},
				fetchNewsError: nil,
			},
			&spfClientMock{
				news: []News{
					{
						Title:   "Cool space fact",
						Summary: "There are spaceflights",
					},
				},
				fetchNewsError: nil,
			},
			"",
			-5,
			http.StatusOK,
			[]News{
				{Title: "Cool space fact", Summary: "There are spaceflights"},
				{Title: "Cat fact 1", Summary: "Cats are cool"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/news?tag="+tt.tag, nil)
			newsHandler := handleNews(tt.catClient, tt.spfClient)
			newsHandler(rr, req)
			response := rr.Result()
			body, _ := ioutil.ReadAll(response.Body)
			response.Body.Close()

			if response.StatusCode != tt.expectedStatus {
				t.Errorf("%s, expected %d, got %d", tt.name, tt.expectedStatus, response.StatusCode)
			}
			var responseData []News
			_ = json.Unmarshal(body, &responseData)
			if !reflect.DeepEqual(responseData, tt.expectedNews) {
				t.Errorf("%s, Response data is not expected", tt.name)
			}
		})
	}
}
