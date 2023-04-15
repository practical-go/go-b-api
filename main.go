package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
	})
	http.HandleFunc("/news", handleNews)
	http.ListenAndServe(":8080", nil)
}

func handleNews(w http.ResponseWriter, r *http.Request) {
	spaceflightNews, err := fetchSpaceflightNews()
	if err != nil {
		http.Error(w, "Error fetching news", http.StatusInternalServerError)
		return
	}
	catFacts, err := fetchCatFacts()
	if err != nil {
		http.Error(w, "Error fetching news", http.StatusInternalServerError)
		return
	}
	var news []News

	for i, sf, cf := 1, 0, 0; i <= 10; i++ {
		if i%3 != 0 && sf < len(spaceflightNews) {
			news = append(news, News{
				Title:   spaceflightNews[sf].Title,
				Summary: spaceflightNews[sf].Summary,
			})
			sf++
		} else if cf < len(catFacts) {
			news = append(news, News{
				Title:   "Cat fact",
				Summary: catFacts[cf].Text,
			})
			cf++
		}
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

type CatFact struct {
	Text string `json:"text"`
}

type SpaceflightNews struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

type News struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

func fetchCatFacts() ([]CatFact, error) {
	body, err := getRequest("https://cat-fact.herokuapp.com/facts/")
	if err != nil {
		return nil, err
	}

	var catFacts []CatFact
	err = json.Unmarshal(body, &catFacts)
	if err != nil {
		return nil, err
	}

	return catFacts, nil
}

func fetchSpaceflightNews() ([]SpaceflightNews, error) {
	body, err := getRequest("https://api.spaceflightnewsapi.net/v4/articles/?limit=10")
	if err != nil {
		return nil, err
	}

	var news []SpaceflightNews
	err = json.Unmarshal(body, &news)
	if err != nil {
		return nil, err
	}

	return news, nil
}

func getRequest(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
