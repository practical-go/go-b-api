package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type newsFetcher interface {
	fetchNews(limit int) ([]News, error)
}

func handleNews(catFactsClient newsFetcher, spfNewsClient newsFetcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tag := r.URL.Query().Get("tag")
		var spaceflightNews, catFacts []News
		limit := getLimit(r.URL.Query().Get("limit"))

		if tag != "cat" {
			var err error
			spaceflightNews, err = spfNewsClient.fetchNews(limit)
			if err != nil {
				http.Error(w, "Error fetching news", http.StatusInternalServerError)
				return
			}
		}
		if tag != "space" {
			var err error
			catFacts, err = catFactsClient.fetchNews(limit)
			if err != nil {
				http.Error(w, "Error fetching news", http.StatusInternalServerError)
				return
			}
		}

		var news []News
		switch tag {
		case "cat":
			news = catFacts
		case "space":
			news = spaceflightNews
		default:
			news = createFeed(limit, spaceflightNews, catFacts)
		}

		jsonData, err := json.Marshal(news)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Write(jsonData)
	}
}

func getLimit(limitStr string) int {
	if limitStr == "" {
		return 10
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit > 50 || limit < 1 {
		return 10
	}
	return limit
}

func createFeed(limit int, spaceflightNews, catFacts []News) []News {
	var news []News
	for i, sf, cf := 1, 0, 0; i <= limit; i++ {
		if i%3 != 0 && sf < len(spaceflightNews) {
			news = append(news, spaceflightNews[sf])
			sf++
		} else if cf < len(catFacts) {
			news = append(news, catFacts[cf])
			cf++
		}
	}
	return news
}
