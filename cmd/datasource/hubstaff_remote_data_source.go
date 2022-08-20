package datasource

import (
	"kitty-bot/cmd/domain/hubstaff"
	"kitty-bot/internal/helper"
	"log"

	"github.com/go-resty/resty/v2"
)

type hubstaffRemoteDataSource struct {
	requestHubstaff     *resty.Request
	requestHubstaffAuth *resty.Request
}

func NewHubstaffRemoteDataSource(
	requestHubstaff *resty.Request,
	requestHubstaffAuth *resty.Request,
) *hubstaffRemoteDataSource {
	return &hubstaffRemoteDataSource{
		requestHubstaff:     requestHubstaff,
		requestHubstaffAuth: requestHubstaffAuth,
	}
}

func (r hubstaffRemoteDataSource) Login() (loginResponse hubstaff.LoginResponse) {
	path := "access_tokens"
	_, err := r.requestHubstaffAuth.
		SetResult(&loginResponse).
		Post(path)
	if err != nil {
		log.Fatal(helper.SprintLog("login hubstaff failed"))
	}
	return
}

func (r hubstaffRemoteDataSource) GetListMembers() (membeResponse hubstaff.MemberResponse) {
	path := "projects/671717/members"
	_, err := r.requestHubstaff.
		SetResult(&membeResponse).
		Get(path)
	if err != nil {
		log.Fatal(helper.SprintLog("get list members failed"))
	}
	return
}

func (r hubstaffRemoteDataSource) GetDailyActivityByRangeDate(
	startDate string,
	stopDate string,
) (dailyActivityResponse hubstaff.DailyActivityResponse) {
	path := "organizations/166190/activities/daily"
	_, err := r.requestHubstaff.
		SetQueryParam("date[start]", startDate).
		SetQueryParam("date[stop]", stopDate).
		SetResult(&dailyActivityResponse).
		Get(path)
	if err != nil {
		log.Fatal(helper.SprintLog("get daily activity by range date failed"))
	}
	return
}
