package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type newsDeleter interface {
	deleteNews(ctx context.Context, uuid string) error
}

func handleDeleteNews(store newsDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "404 Does not exist", http.StatusNotFound)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		uuid := r.URL.Query().Get("uuid")

		if uuid != "" {
			err := store.deleteNews(ctx, uuid)
			if err != nil {
				http.Error(w, "Error deleting news", http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "UUID not provided", http.StatusBadRequest)
			return
		}

		fmt.Fprint(w, "Successfully deleted")
	}
}
