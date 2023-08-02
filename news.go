package main

type News struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

type NewsExtended struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Id      string
}
