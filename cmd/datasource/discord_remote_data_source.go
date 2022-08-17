package datasource

import (
	"fmt"
	"kitty-bot/cmd/domain/discord"
	"kitty-bot/configs"
	"kitty-bot/internal/helper"
	"log"
	"strings"

	"github.com/go-resty/resty/v2"
)

type discordRemoteDataSource struct {
	request *resty.Request
	baseApi string
}

func NewDiscordRemoteDataSource(request *resty.Request) *discordRemoteDataSource {
	baseApi := helper.LoadEnvVariable(configs.BaseApiDiscord)
	return &discordRemoteDataSource{request: request, baseApi: baseApi}
}

func (r *discordRemoteDataSource) SendMessageToChannel(idChannel string, body discord.DiscordCreateMessageBody) (result bool) {
	path := r.baseApi + "channels/" + idChannel + "/messages"
	discordBotToken := helper.LoadEnvVariable(configs.DiscordBotToken)
	response, err := r.request.
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bot %s", discordBotToken)).
		SetBody(body).
		Post(path)
	if err != nil {
		log.Fatal(helper.SprintLog("send message to channel failed"))
	}
	statusCode := response.StatusCode()
	strStatusCode := fmt.Sprint(statusCode)
	if strings.HasPrefix(strStatusCode, "2") {
		result = true
	} else {
		result = false
	}
	return
}
