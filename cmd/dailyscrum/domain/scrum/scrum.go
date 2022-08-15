package scrum

type Scrum struct {
	Data []ItemScrum `json:"data"`
}

type ItemScrum struct {
	Message string `json:"message"`
	Image   string `json:"image"`
}
