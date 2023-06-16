package main

import (
	"fmt"
	"net/http"
)

func main() {
	var catFactsClient CatfactClient
	catFactsClient.Init()
	var spfNewsClient SpaceflightClient
	spfNewsClient.Init()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
	})
	http.HandleFunc("/news", handleNews(&catFactsClient, &spfNewsClient))
	http.ListenAndServe(":8080", nil)
}
