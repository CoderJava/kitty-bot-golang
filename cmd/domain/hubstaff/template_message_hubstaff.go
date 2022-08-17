package hubstaff

type TemplateMessageHubstaff struct {
	IdHubstaff string `json:"id_hubstaff"`
	IdDiscord  string `json:"id_discord"`
	Name       string `json:"name"`
	Tracked    int    `json:"tracked"`
	Idle       int    `json:"idle"`
}
