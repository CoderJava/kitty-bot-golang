package datasource

import (
	"fmt"
	"kitty-bot/cmd/domain/discord"
	"kitty-bot/internal/helper"
	"log"
	"strings"

	"github.com/go-resty/resty/v2"
)

type discordRemoteDataSource struct {
	requestDiscord *resty.Request
}

func NewDiscordRemoteDataSource(requestDiscord *resty.Request) *discordRemoteDataSource {
	return &discordRemoteDataSource{requestDiscord: requestDiscord}
}

func (r *discordRemoteDataSource) SendMessageToChannel(idChannel string, body discord.DiscordCreateMessageBody) (result bool) {
	path := "channels/" + idChannel + "/messages"
	response, err := r.requestDiscord.
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
