// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Task struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Entry       string   `json:"entry"`
	Modified    string   `json:"modified"`
	UUID        string   `json:"uuid"`
	Urgency     float64  `json:"urgency"`
	Status      string   `json:"status"`
	Priority    string   `json:"priority"`
	Due         string   `json:"due"`
	Project     string   `json:"project"`
	Tags        []string `json:"tags"`
}

type TimeRecord struct {
	ID    string   `json:"id"`
	Start string   `json:"start"`
	End   string   `json:"end"`
	Tags  []string `json:"tags"`
}
