package hubstaff

type DailyActivityResponse struct {
	DailyActivities []ItemDailyActivityResponse `json:"daily_activities"`
	Pagination PaginationDailyActivityResponse `json:"pagination"`
}

type ItemDailyActivityResponse struct {
	Date                  string `json:"date"`
	UserId                int    `json:"user_id"`
	Overall               int    `json:"overall"`
	InputTrackedInSeconds int    `json:"input_tracked"`
	IdleInSeconds         int    `json:"idle"`
}

type PaginationDailyActivityResponse struct {
	NextPageStartId int `json:"next_page_start_id"`
}