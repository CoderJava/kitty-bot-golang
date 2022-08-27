package app

import (
	"encoding/json"
	"fmt"
	"kitty-bot/cmd/datasource"
	"kitty-bot/cmd/domain/discord"
	"kitty-bot/cmd/domain/event"
	"kitty-bot/cmd/domain/hubstaff"
	"kitty-bot/configs"
	"kitty-bot/internal/helper"
	"log"
	"os"
	"sort"
	"time"

	"github.com/go-resty/resty/v2"
)

func StartDailyHubstaff(
	requestHubstaff *resty.Request,
	requestHubstaffAuth *resty.Request,
	requestDiscord *resty.Request,
	cacheHelper helper.CacheHelper,
) {
	// baca APP_ENV
	appEnv := helper.LoadEnvVariable(configs.AppEnv)

	// hit ke endpoint login hubstaff dan simpan access tokennya didalam cache
	hubstaffRemoteDataSource := datasource.NewHubstaffRemoteDataSource(requestHubstaff, requestHubstaffAuth)
	loginResponse := hubstaffRemoteDataSource.Login()
	if (hubstaff.LoginResponse{}) == loginResponse {
		helper.PrintLog("endpoint login hubstaff gagal")
		return
	}
	cacheHelper.Set(configs.AccessToken, loginResponse.AccessToken)
	cacheHelper.Set(configs.RefreshToken, loginResponse.RefreshToken)

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

	patternDate := "2006-01-02"
	strStartDay := startDate.Format(patternDay)
	strStartDate := startDate.Format(patternDate)
	strStopDate := stopDate.Format(patternDate)

	// cek apakah [strStartDate] merupakan hari libur nasional
	for _, itemEventHoliday := range listEventHoliday {
		if strStartDate == itemEventHoliday.StrDate && itemEventHoliday.IsHoliday {
			helper.PrintLog("daily hubstaff hari ini adalah hari libur")
			return
		}
	}

	// tentukan syarat total jam kerja
	var requirementWorkingHourInSeconds int
	getRequirementWorkingHourInSeconds(strStartDay, &requirementWorkingHourInSeconds)
	if requirementWorkingHourInSeconds == 0 {
		helper.PrintLog("daily hubstaff requirement working hour is zero")
		return
	}

	dailyActivityResponse := hubstaffRemoteDataSource.GetDailyActivityByRangeDate(strStartDate, strStopDate, 0)
	listMessages := []hubstaff.TemplateMessageHubstaff{}
	for _, itemDailyActivityResponse := range dailyActivityResponse.DailyActivities {
		userIdHubstaff := fmt.Sprint(itemDailyActivityResponse.UserId)

		// user hubstaff yang tidak dihitung
		if userIdHubstaff == helper.LoadEnvVariable(configs.IdHubstaffAdeIskandar) {
			continue
		}

		inputTrackedInSeconds := itemDailyActivityResponse.InputTrackedInSeconds
		idleInSeconds := itemDailyActivityResponse.IdleInSeconds

		// jika id discord-nya tidak valid
		name, idDiscord := helper.GetNameAndIdDiscordByIdHubstaff(userIdHubstaff)
		if idDiscord == "" {
			continue
		}

		listAlready := helper.FilterTemplateMessageHubstaff(listMessages, func(index int) bool {
			return listMessages[index].IdHubstaff == userIdHubstaff
		})
		if len(listAlready) > 0 {
			// untuk menghitung jika lebih dari satu record daily activity hubstaff-nya
			for _, elementAlready := range listAlready {
				inputTrackedInSeconds += elementAlready.Tracked
				idleInSeconds += elementAlready.Idle
			}
			indexUpdate := helper.SliceIndex(len(listMessages), func(index int) bool {
				return listMessages[index].IdHubstaff == userIdHubstaff
			})
			listMessages[indexUpdate] = hubstaff.TemplateMessageHubstaff{
				IdHubstaff: userIdHubstaff,
				IdDiscord:  idDiscord,
				Name:       name,
				Tracked:    inputTrackedInSeconds,
				Idle:       idleInSeconds,
			}
		} else {
			listMessages = append(listMessages, hubstaff.TemplateMessageHubstaff{
				IdHubstaff: userIdHubstaff,
				IdDiscord:  idDiscord,
				Name:       name,
				Tracked:    inputTrackedInSeconds,
				Idle:       idleInSeconds,
			})
		}
	}

	// jika tidak ada user yang bekerja pada hari tersebut maka, jangan kirimkan info daily hubstaff-nya
	if len(listMessages) == 0 {
		helper.PrintLog("list message hubstaff kosong")
		return
	}

	// filter user yang jam kerjanya belum terpenuhi
	listMessages = helper.FilterTemplateMessageHubstaff(listMessages, func(index int) bool {
		element := listMessages[index]
		return element.Tracked < requirementWorkingHourInSeconds
	})

	// sort dari yang paling besar ke terkecil jam hubstaff-nya
	sort.Slice(listMessages, func(i, j int) bool {
		return listMessages[i].Tracked > listMessages[j].Tracked
	})

	strDate := startDate.Format("02-01-2006")
	strRequirementWorkingHour := helper.ConvertSecondToFormatHourMinuteSecond(requirementWorkingHourInSeconds)
	content := "@everyone Ladies and gentlemen..."
	content += "\n:trophy: Welcome to the HubStaff Championship :trophy:"
	content += fmt.Sprintf(
		"\n\nBerikut adalah para peserta yang kalah pada tanggal **%s** dengan syarat jam HubStaff-nya **%s**.",
		strDate,
		strRequirementWorkingHour,
	)
	content += "\n"

	for index, itemMessage := range listMessages {
		content += fmt.Sprintf("%d. <@%s>", index+1, itemMessage.IdDiscord)
		content += fmt.Sprintf("\nTotal\t: **%s**", helper.ConvertSecondToFormatHourMinuteSecond(itemMessage.Tracked))
		content += fmt.Sprintf("\nIdle\t** ** : **%s**", helper.ConvertSecondToFormatHourMinuteSecond(itemMessage.Idle))
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
	isSuccessSendMessageDailyHubstaff := discordRemoteDataSource.SendMessageToChannel(
		idChannel,
		discord.DiscordCreateMessageBody{
			Content: content,
			Embeds:  []discord.ItemEmbedDiscord{},
		})
	if isSuccessSendMessageDailyHubstaff {
		helper.PrintLog("daily hubstaff success")
	} else {
		helper.PrintLog("daily hubstaff failure")
	}
}

func getRequirementWorkingHourInSeconds(strDay string, second *int) {
	switch strDay {
	case "Mon", "Tue", "Wed", "Thu":
		// 07:30
		*second = (3600 * 7) + (60 * 30)
	case "Fri":
		// 06:30
		*second = (3600 * 6) + (60 * 30)
	case "Sat":
		// 03:30
		*second = (3600 * 3) + (60 * 30)
	}
}
