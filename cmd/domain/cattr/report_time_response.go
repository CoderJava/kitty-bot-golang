package cattr

type ReportTimeResponse struct {
	Data []DataReportTimeResponse `json:"data"`
}

type DataReportTimeResponse struct {
	TimeInSeconds int                    `json:"time"`
	User          UserReportTimeResponse `json:"user"`
}

type UserReportTimeResponse struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Fullname string `json:"full_name"`
}
