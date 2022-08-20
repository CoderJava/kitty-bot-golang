package helper

import (
	"fmt"
	"kitty-bot/configs"

	"github.com/go-resty/resty/v2"
)

type middlewareHelper struct {
	cacheHelper CacheHelper
}

func NewMiddlewareHelper(cacheHelper CacheHelper) *middlewareHelper {
	return &middlewareHelper{cacheHelper: cacheHelper}
}

func (mh middlewareHelper) OnBeforeRequestDiscord(c *resty.Client, req *resty.Request) error {
	baseApi := LoadEnvVariable(configs.BaseApiDiscord)
	discordBotToken := LoadEnvVariable(configs.DiscordBotToken)

	c.SetBaseURL(baseApi)
	req.SetHeader("Content-Type", "application/json")
	req.SetHeader("Authorization", fmt.Sprintf("Bot %s", discordBotToken))
	return nil
}

func (mh middlewareHelper) OnBeforeRequestHubstaffAuth(c *resty.Client, req *resty.Request) error {
	baseApi := LoadEnvVariable(configs.BaseApiAuthHubstaff)
	refreshToken := mh.cacheHelper.Get(configs.RefreshToken)
	if refreshToken == "" {
		refreshToken = LoadEnvVariable(configs.RefreshToken)
	}

	c.SetBaseURL(baseApi)
	req.SetHeader("Content-Type", "application/json")
	req.SetQueryParam("grant_type", "refresh_token")
	req.SetQueryParam("refresh_token", refreshToken)
	return nil
}

func (mh middlewareHelper) OnBeforeRequestHubstaff(c *resty.Client, req *resty.Request) error {
	baseApi := LoadEnvVariable(configs.BaseApiHubstaff)
	accessToken := mh.cacheHelper.Get(configs.AccessToken)

	c.SetBaseURL(baseApi)
	req.SetHeader("Content-Type", "application/json")
	req.SetHeader("Authorization", "Bearer "+accessToken)
	return nil
}