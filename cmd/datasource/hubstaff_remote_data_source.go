package datasource

import (
	"kitty-bot/cmd/domain/hubstaff"
	"kitty-bot/configs"
	"kitty-bot/internal/helper"
	"log"

	"github.com/go-resty/resty/v2"
)

type hubstaffRemoteDataSource struct {
	request     *resty.Request
	baseApi     string
	baseApiAuth string
	cacheHelper helper.CacheHelper
}

func NewHubstaffRemoteDataSource(
	request *resty.Request,
	cacheHelper helper.CacheHelper,
) *hubstaffRemoteDataSource {
	baseApi := helper.LoadEnvVariable(configs.BaseApiHubstaff)
	baseApiAuth := helper.LoadEnvVariable(configs.BaseApiAuthHubstaff)
	return &hubstaffRemoteDataSource{
		request:     request,
		baseApi:     baseApi,
		baseApiAuth: baseApiAuth,
		cacheHelper: cacheHelper,
	}
}

func (r hubstaffRemoteDataSource) Login() (loginResponse hubstaff.LoginResponse) {
	path := r.baseApiAuth + "access_tokens"
	refreshToken := r.cacheHelper.Get(configs.RefreshToken)
	if refreshToken == "" {
		refreshToken = helper.LoadEnvVariable(configs.RefreshToken)
	}
	_, err := r.request.
		SetHeader("Content-Type", "application/json").
		SetQueryParam("grant_type", "refresh_token").
		SetQueryParam("refresh_token", refreshToken).
		SetResult(&loginResponse).
		Post(path)
	if err != nil {
		log.Fatal(helper.SprintLog("login hubstaff failed"))
	}
	return
}

func (r hubstaffRemoteDataSource) GetListMembers() (membeResponse hubstaff.MemberResponse) {
	path := r.baseApi + "projects/671717/members"
	accessToken := r.cacheHelper.Get(configs.AccessToken)
	_, err := r.request.
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+accessToken).
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
	path := r.baseApi + "organizations/166190/activities/daily"
	accessToken := r.cacheHelper.Get(configs.AccessToken)
	_, err := r.request.
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+accessToken).
		SetQueryParam("date[start]", startDate).
		SetQueryParam("date[stop]", stopDate).
		SetResult(&dailyActivityResponse).
		Get(path)
	if err != nil {
		log.Fatal(helper.SprintLog("get daily activity by range date failed"))
	}
	return
}
