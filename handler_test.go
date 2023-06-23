package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type catClientMock struct {
	fetchNewsError error
	news           []News
}

func (c *catClientMock) fetchNews(limit int) ([]News, error) {
	return nil, c.fetchNewsError
}

type spfClientMock struct {
	fetchNewsError error
	news           []News
}

func (c *spfClientMock) fetchNews(limit int) ([]News, error) {
	return nil, c.fetchNewsError
}

func TestNewsHandler(t *testing.T) {
	tests := []struct {
		name           string
		catClient      newsFetcher
		spfClient      newsFetcher
		tag            string
		limit          int
		expectedStatus int
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
		},
		{
			"Failure",
			&catClientMock{
				news:           nil,
				fetchNewsError: errors.New("Some error"),
			},
			&spfClientMock{
				news:           nil,
				fetchNewsError: errors.New("Another error"),
			},
			"",
			10,
			http.StatusInternalServerError,
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/news?tag="+tt.tag, nil)
			newsHandler := handleNews(tt.catClient, tt.spfClient)
			newsHandler(rr, req)
			if rr.Result().StatusCode != tt.expectedStatus {
				t.Errorf("%s, expected %d, got %d", tt.name, tt.expectedStatus, rr.Result().StatusCode)
			}
		})
	}
}
