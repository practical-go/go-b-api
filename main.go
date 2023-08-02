package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
	})
	http.HandleFunc("/news", handleNews(&catFactsClient, &spfNewsClient))
	http.HandleFunc("/createnews", handleUpsertNews(&dbClient))
	http.HandleFunc("/dbnews", handleDBNews(&dbClient))
	http.HandleFunc("/deletenews", handleDeleteNews(&dbClient))
	http.ListenAndServe(":8080", nil)

}
