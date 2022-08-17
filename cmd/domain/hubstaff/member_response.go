package hubstaff

type MemberResponse struct {
	Members []ItemMemberResponse `json:"members"`
}

type ItemMemberResponse struct {
	UserId int    `json:"user_id"`
	Status string `json:"membership_status"`
}
