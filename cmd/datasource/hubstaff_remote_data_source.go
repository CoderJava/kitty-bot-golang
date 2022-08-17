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
}

func NewHubstaffRemoteDataSource(request *resty.Request) *hubstaffRemoteDataSource {
	baseApi := helper.LoadEnvVariable(configs.BaseApiHubstaff)
	baseApiAuth := helper.LoadEnvVariable(configs.BaseApiAuthHubstaff)
	return &hubstaffRemoteDataSource{
		request:     request,
		baseApi:     baseApi,
		baseApiAuth: baseApiAuth,
	}
}

func (r hubstaffRemoteDataSource) Login() (loginResponse hubstaff.LoginResponse) {
	path := r.baseApiAuth + "access_tokens"
	refreshToken := helper.LoadEnvVariable(configs.RefreshToken)
	_, err := r.request.
		SetHeader("Content-Type", "application/json").
		SetQueryParam("grant_type", "refresh_token").
		SetQueryParam("refresh_token", refreshToken).
		SetResult(loginResponse).
		Post(path)
	if err != nil {
		log.Fatal(helper.SprintLog("login hubstaff failed"))
	}
	return
}

/* func (r hubstaffRemoteDataSource) GetListMembers() (membeResponse hubstaff.MemberResponse) {
	path := r.baseApi + "projects/671717/members"
	_, err := r.request.
		SetHeader("Content-Type", "application/json")

	return
} */
