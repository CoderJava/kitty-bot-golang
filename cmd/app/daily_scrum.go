package app

import (
	"encoding/json"
	"kitty-bot/cmd/datasource"
	"kitty-bot/cmd/domain/discord"
	"kitty-bot/cmd/domain/event"
	"kitty-bot/cmd/domain/scrum"
	"kitty-bot/configs"
	"kitty-bot/internal/helper"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

func StartDailyScrum(requestDiscord *resty.Request, listDaysScrum []string) {
	// baca APP_ENV dari environment variable
	appEnv := helper.LoadEnvVariable(configs.AppEnv)

	// baca file template_scrum.json
	// dan ubah menjadi struct scrum
	templateScrumJson, err := os.ReadFile(configs.TemplateMessageDailyScrum)
	if err != nil {
		log.Fatal(err, helper.SprintLog("gagal baca file template scrum"))
	}
	var templateMessageScrum scrum.TemplateMessageScrum
	err = json.Unmarshal(templateScrumJson, &templateMessageScrum)
	if err != nil {
		log.Fatal(err, helper.SprintLog("gagal parsing json string kedalam struct template message scrum"))
	}
	listTemplateMessages := templateMessageScrum.Data

	// baca file event_holiday.json
	// dan ubah menjadi struct event
	eventJson, err := os.ReadFile(configs.EventHoliday)
	if err != nil {
		log.Fatal(err, helper.SprintLog("gagal baca file event holiday"))
	}
	var listEvents event.Event
	err = json.Unmarshal(eventJson, &listEvents)
	if err != nil {
		log.Fatal(err, helper.SprintLog("gagal parsing json string kedalam struct event"))
	}

	now := time.Now()
	formattedDay := now.Format("Mon")
	formattedDate := now.Format("2006-01-02")

	// cek apakah hari ini ada jadwal daily scrum
	isDailyScrum := helper.ContainString(listDaysScrum, formattedDay)
	if !isDailyScrum {
		helper.PrintLog("Hari ini tidak ada jadwal daily scrum")
		return
	}

	// cek apakah hari ini libur nasional
	var isHoliday bool
	for _, itemEvent := range listEvents.Data {
		if itemEvent.StrDate == formattedDate && itemEvent.IsHoliday {
			isHoliday = true
			break
		}
	}
	if isHoliday {
		helper.PrintLog("Hari ini adalah hari libur nasional")
		return
	}

	// ambil secara acak message reminder scrum
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(listTemplateMessages))
	selectedTemplateMessage := listTemplateMessages[randomIndex]

	// Kirim daily reminder scrum ke channel discord
	discordCreateMessageBody := discord.DiscordCreateMessageBody{
		Content: "",
		Embeds: []discord.ItemEmbedDiscord{
			{
				Title:       "Scrum Reminder",
				Description: selectedTemplateMessage.Message,
				Image: discord.ImageItemEmbedDiscord{
					Url: selectedTemplateMessage.Image,
				},
			},
		},
	}
	discordRemoteDataSource := datasource.NewDiscordRemoteDataSource(requestDiscord)
	var idChannelDiscord string
	if appEnv == "development" {
		idChannelDiscord = helper.LoadEnvVariable(configs.IdChannelDiscordDevelopment)
	} else if appEnv == "production" {
		idChannelDiscord = helper.LoadEnvVariable(configs.IdChannelDiscordProduction)
	}
	isSuccessSendMessageDailyScrum := discordRemoteDataSource.SendMessageToChannel(
		idChannelDiscord,
		discordCreateMessageBody,
	)
	if isSuccessSendMessageDailyScrum {
		helper.PrintLog("Success send message daily scrum")
	} else {
		helper.PrintLog("Failure send message daily scrum")
	}
}
