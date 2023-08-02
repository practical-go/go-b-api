package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	var catFactsClient CatfactClient
	catFactsClient.Init()
	var spfNewsClient SpaceflightClient
	spfNewsClient.Init()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var dbClient dbClient
	dbClient.Init(ctx, "newsdb", "news")
	defer func() { dbClient.Close(ctx) }()

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
	}).Methods("GET")
	router.HandleFunc("/news", handleNews(&catFactsClient, &spfNewsClient)).Methods("GET")
	router.HandleFunc("/dbnews", handleUpsertNews(&dbClient)).Methods("POST")
	router.HandleFunc("/dbnews", handleDBNews(&dbClient)).Methods("GET")
	router.HandleFunc("/dbnews", handleDeleteNews(&dbClient)).Methods("DELETE")

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)

}
