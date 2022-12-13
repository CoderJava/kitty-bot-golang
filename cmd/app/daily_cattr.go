package app

import (
	"encoding/json"
	"fmt"
	"kitty-bot/cmd/datasource"
	"kitty-bot/cmd/domain/cattr"
	"kitty-bot/cmd/domain/discord"
	"kitty-bot/cmd/domain/event"
	"kitty-bot/configs"
	"kitty-bot/internal/helper"
	"log"
	"os"
	"sort"
	"time"

	"github.com/go-resty/resty/v2"
)

func StartDailyCattr(
	requestCattr *resty.Request,
	requestDiscord *resty.Request,
	cacheHelper helper.CacheHelper,
) {
	// baca APP_ENV
	appEnv := helper.LoadEnvVariable(configs.AppEnv)

	// hit ke endpoint login cattr dan simpan access tokennya didalam cache
	cattrRemoteDataSource := datasource.NewCattrRemoteDataSource(requestCattr)
	usernameCattr := helper.LoadEnvVariable(configs.UsernameCattr)
	passwordCattr := helper.LoadEnvVariable(configs.PasswordCattr)
	loginBody := cattr.LoginBody{
		Email:    usernameCattr,
		Password: passwordCattr,
	}
	loginResponse := cattrRemoteDataSource.Login(loginBody)
	if (cattr.LoginResponse{}) == loginResponse {
		helper.PrintLog("endpoint login cattr gagal")
		return
	}
	cacheHelper.Set(configs.AccessToken, loginResponse.Data.AccessToken)

	// pastikan hari ini adalah hari kerja
	now := time.Now()
	patternDay := "Mon"
	formattedDay := now.Format(patternDay)
	listDays := []string{
		"Mon",
		"Tue",
		"Wed",
		"Thu",
		"Fri",
		"Sat",
	}
	isTodayWorking := helper.ContainString(listDays, formattedDay)
	if !isTodayWorking {
		return
	}

	// ambil data json event yang di asset
	jsonEventHoliday, err := os.ReadFile(configs.EventHoliday)
	if err != nil {
		log.Fatal(err, helper.SprintLog("gagal baca file event holiday"))
	}
	var listEvents event.Event
	err = json.Unmarshal(jsonEventHoliday, &listEvents)
	if err != nil {
		log.Fatal(err, helper.SprintLog("gagal parsing json event holiday"))
	}
	listEventHoliday := listEvents.Data

	var startDate time.Time
	var stopDate time.Time

	if formattedDay == "Mon" {
		// khusus untuk hari senin ambil ke hari sabtu yang mana itu H-2
		tempDate := now.AddDate(0, 0, -2)
		startDate = tempDate
		stopDate = tempDate
	} else {
		// selain hari senin maka, ambil ke hari sebelumnya H-1
		tempDate := now.AddDate(0, 0, -1)
		startDate = tempDate
		stopDate = tempDate
	}

	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	stopDate = time.Date(stopDate.Year(), stopDate.Month(), stopDate.Day(), 23, 59, 59, 0, stopDate.Location())

	patternDate := "2006-01-02T15:04:05Z0700"
	strStartDay := startDate.Format(patternDay)
	strStartDate := startDate.Format(patternDate)
	strStopDate := stopDate.Format(patternDate)

	// cek apakah [strStartDate] merupakan hari libur nasional
	for _, itemEventHoliday := range listEventHoliday {
		if strStartDate == itemEventHoliday.StrDate && itemEventHoliday.IsHoliday {
			helper.PrintLog("daily cattr hari ini adalah hari libur")
			return
		}
	}

	// tentukan syarat total jam kerja
	var requirementWorkingHourInSeconds int
	helper.GetRequirementWorkingHourInSeconds(strStartDay, &requirementWorkingHourInSeconds)
	if requirementWorkingHourInSeconds == 0 {
		helper.PrintLog("daily cattr requirement working hour is zero")
		return
	}

	users := []string{}
	// users = append(users, helper.LoadEnvVariable(configs.IdCattrAdeIskandar)) // user cattr yang tidak dihitung
	users = append(users, helper.LoadEnvVariable(configs.IdCattrYudiSetiawan))
	users = append(users, helper.LoadEnvVariable(configs.IdCattrRyan))
	users = append(users, helper.LoadEnvVariable(configs.IdCattrSabrino))
	users = append(users, helper.LoadEnvVariable(configs.IdCattrRioDwi))
	users = append(users, helper.LoadEnvVariable(configs.IdCattrBobby))
	users = append(users, helper.LoadEnvVariable(configs.IdCattrAditama))
	users = append(users, helper.LoadEnvVariable(configs.IdCattrAldoFaiz))
	users = append(users, helper.LoadEnvVariable(configs.IdCattrDewi))
	users = append(users, helper.LoadEnvVariable(configs.IdCattrAbdulAziz))
	users = append(users, helper.LoadEnvVariable(configs.IdCattrRianto))
	users = append(users, helper.LoadEnvVariable(configs.IdCattrAbdi))
	reportTimeBody := cattr.ReportTimeBody{
		Users:   users,
		StartAt: strStartDate,
		EndAt:   strStopDate,
	}

	reportTimeResponse := cattrRemoteDataSource.GetReportTime(reportTimeBody)
	listMessages := []cattr.TemplateMessageCattr{}
	for _, itemReportTimeResponse := range reportTimeResponse.Data {
		userIdCattr := fmt.Sprint(itemReportTimeResponse.User.Id)

		// jika id discord-nya tidak valid
		name, idDiscord := helper.GetNameAndIdDiscordByIdCattr(userIdCattr)
		if idDiscord == "" {
			continue
		}

		totalTrackedInSeconds := itemReportTimeResponse.TimeInSeconds
		listMessages = append(listMessages, cattr.TemplateMessageCattr{
			IdCattr:   userIdCattr,
			IdDiscord: idDiscord,
			Name:      name,
			Tracked:   totalTrackedInSeconds,
		})
	}

	// jika tidak ada user yang bekerja pada hari tersebut maka, jangan kirimkan info daily cattr-nya
	if len(listMessages) == 0 {
		helper.PrintLog("list message cattr kosong")
		return
	}

	// filter user yang jam kerjanya belum terpenuhi
	listMessages = helper.FilterTemplateMessageCattr(listMessages, func(index int) bool {
		element := listMessages[index]
		return element.Tracked < requirementWorkingHourInSeconds
	})

	// sort dari yang paling besar ke terkecil jam cattr-nya
	sort.Slice(listMessages, func(i, j int) bool {
		return listMessages[i].Tracked > listMessages[j].Tracked
	})

	strDate := startDate.Format("02-01-2006")
	strRequirementWorkingHour := helper.ConvertSecondToFormatHourMinuteSecond(requirementWorkingHourInSeconds)
	content := "@everyone Ladies and gentlemen..."
	content += "\n:trophy: Welcome to the Cattr Championship :trophy:"
	content += fmt.Sprintf(
		"\n\nBerikut adalah para peserta yang kalah pada tanggal **%s** dengan syarat jam kerjanya **%s**.",
		strDate,
		strRequirementWorkingHour,
	)
	content += "\n"

	for index, itemMessage := range listMessages {
		content += fmt.Sprintf("%d. <@%s>", index+1, itemMessage.IdDiscord)
		content += fmt.Sprintf("\nTotal\t: **%s**", helper.ConvertSecondToFormatHourMinuteSecond(itemMessage.Tracked))
		if index != len(listMessages)-1 {
			content += "\n\n"
		}
	}

	if appEnv == "development" {
		helper.PrintLog(fmt.Sprintf("content: %s", content))
	}

	discordRemoteDataSource := datasource.NewDiscordRemoteDataSource(requestDiscord)
	var idChannel string
	if appEnv == "development" {
		idChannel = helper.LoadEnvVariable(configs.IdChannelDiscordDevelopment)
	} else if appEnv == "production" {
		idChannel = helper.LoadEnvVariable(configs.IdChannelDiscordProduction)
	}
	isSuccessSendMessageDailyCattr := discordRemoteDataSource.SendMessageToChannel(
		idChannel,
		discord.DiscordCreateMessageBody{
			Content: content,
			Embeds:  []discord.ItemEmbedDiscord{},
		},
	)
	if isSuccessSendMessageDailyCattr {
		helper.PrintLog("daily cattr success")
	} else {
		helper.PrintLog("daily cattr failure")
	}
}
