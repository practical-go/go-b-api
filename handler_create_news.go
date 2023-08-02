package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type newsUpserter interface {
	upsertNews(ctx context.Context, title, summary, uuid string) error
}

func handleUpsertNews(store newsUpserter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "404 Does not exist", http.StatusNotFound)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var uuid string
		uuid = r.URL.Query().Get("uuid")
		if uuid == "" {
			uuid = createUUID()
		}
		err := store.upsertNews(ctx, "Japanese Zoo update", "Click to read more!", uuid)
		if err != nil {
			http.Error(w, "Error upserting news", http.StatusInternalServerError)
			return
		}

		//w.Header().Set("Content-Type", "application/json")
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		//w.Header().Set("Access-Control-Allow-Methods", "POST")
		fmt.Fprint(w, uuid)

	}
}

func createUUID() string {
	id := uuid.New()
	return id.String()
}
