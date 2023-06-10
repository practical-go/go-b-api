package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type catClientMock struct {
	fetchNewsError error
}

func (c *catClientMock) fetchNews(limit int) ([]News, error) {
	return nil, c.fetchNewsError
}

type spfClientMock struct {
	fetchNewsError error
}

func (c *spfClientMock) fetchNews(limit int) ([]News, error) {
	return nil, c.fetchNewsError
}

func TestNewsHandler(t *testing.T) {
	tests := []struct {
		name      string
		catClient newsFetcher
		spfClient newsFetcher
		tag       string
		status    int
	}{
		{
			"Failure",
			&catClientMock{errors.New("Some error")},
			&spfClientMock{errors.New("Another error")},
			"",
			http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			newsHandler := handleNews(tt.catClient, tt.spfClient)
			result := newsHandler()
			if result != tt.status {
				t.Errorf("%s, expected %d, got %d", tt.name, tt.status, result)
			}
		})
	}
}
