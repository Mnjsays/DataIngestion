package models

type Posts struct {
	Data       []Source
	IngestedAt string `json:"ingested_at"`
	Source     string `json:"source"`
}
type Source struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}
