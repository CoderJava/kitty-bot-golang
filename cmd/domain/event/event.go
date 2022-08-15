package event

type Event struct {
	Data []ItemEvent `json:"data"`
}

type ItemEvent struct {
	StrDate   string `json:"date"`
	IsHoliday bool   `json:"is_holiday"`
}
