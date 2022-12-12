package datasource

import (
	"kitty-bot/cmd/domain/cattr"
	"kitty-bot/internal/helper"
	"log"

	"github.com/go-resty/resty/v2"
)

type cattrRemoteDataSource struct {
	request *resty.Request
}

func NewCattrRemoteDataSource(request *resty.Request) *cattrRemoteDataSource {
	return &cattrRemoteDataSource{request: request}
}

func (r cattrRemoteDataSource) Login(loginBody cattr.LoginBody) (loginRespnse cattr.LoginResponse) {
	path := "auth/login"
	_, err := r.request.
		SetBody(loginBody).
		SetResult(&loginRespnse).
		Post(path)
	if err != nil {
		log.Fatal(helper.SprintLog("login cattr failed"))
	}
	return
}

func (r cattrRemoteDataSource) GetReportTime(reportTimeBody cattr.ReportTimeBody) (reportTimeResponse cattr.ReportTimeResponse) {
	path := "report/time"
	_, err := r.request.
		SetBody(reportTimeBody).
		SetResult(&reportTimeResponse).
		Post(path)
	if err != nil {
		log.Fatal(helper.SprintLog("get report time failed"))
	}
	return
}
