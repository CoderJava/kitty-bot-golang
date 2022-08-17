package hubstaff

type DailyActivityResponse struct {
	DailyActivities []ItemDailyActivityResponse `json:"daily_activities"`
}

type ItemDailyActivityResponse struct {
	Date                  string `json:"date"`
	UserId                int    `json:"user_id"`
	Overall               int    `json:"overall"`
	InputTrackedInSeconds int    `json:"input_tracked"`
	IdleInSeconds         int    `json:"idle"`
}