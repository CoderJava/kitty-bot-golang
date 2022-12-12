package cattr

type ReportTimeBody struct {
	Users   []string `json:"users"`
	StartAt string   `json:"start_at"`
	EndAt   string   `json:"end_at"`
}
