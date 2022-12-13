package app

import (
	"fmt"
	"kitty-bot/cmd/datasource"
	"kitty-bot/cmd/domain/cattr"
	"kitty-bot/cmd/domain/discord"
	"kitty-bot/configs"
	"kitty-bot/internal/helper"
	"sort"
	"time"

	"github.com/go-resty/resty/v2"
)

func StartMonthlyCattr(
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

	// ambil waktu saat ini
	// dan tentukan tanggal mulai dan akhir
	now := time.Now()
	var startDate time.Time
	var stopDate time.Time
	helper.GetPeriod(now, &startDate, &stopDate)

	// formatting variable startDate dan stopDate
	patternDate := "2006-01-02T15:04:05Z0700"
	strStartDate := startDate.Format(patternDate)
	strStopDate := stopDate.Format(patternDate)

	users := []string{}
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

		name, idDiscord := helper.GetNameAndIdDiscordByIdCattr(userIdCattr)

		// abaikan datanya jika id discord tidak valid
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

	listBestTracked := []cattr.TemplateMessageCattr{}
	listBestTracked = append(listBestTracked, listMessages...)

	// urutkan user yang best tracked secara ascending
	sort.Slice(listBestTracked, func(x, y int) bool {
		totalX := listBestTracked[x].Tracked
		totalY := listBestTracked[y].Tracked
		return totalX < totalY
	})

	strNowYear := startDate.Format("2006")
	patternDatePeriod := "02 Jan"
	strPeriod := startDate.Format(patternDatePeriod) + " - " + stopDate.Format(patternDatePeriod)
	content := "@everyone Welcome ladies and gentlemen to the"
	content += "\n:trophy: Cattr Championship " + strNowYear + " :trophy:"
	content += "\n\n"
	content += "Setelah sebulan lamanya para peserta kita bekerja dengan **serius tapi santai** untuk memenangkan kompetisi ini "
	content += "maka, tibalah waktunya untuk mengumumkan the next winner of Cattr Championship " + strNowYear + " periode **" + strPeriod + "**."

	// buat pesan untuk tracked terbaik
	content += "\n\n"
	content += ":timer: **Tracked Terbaik**"
	content += "\n================="
	counter := 0
	for index := len(listBestTracked) - 1; index >= 0; index-- {
		counter += 1
		element := listBestTracked[index]
		var strEmojiMedal string
		switch counter {
		case 1:
			strEmojiMedal = ":first_place:"
		case 2:
			strEmojiMedal = ":second_place:"
		case 3:
			strEmojiMedal = ":third_place:"
		}
		content += fmt.Sprintf("\n%s <@%s>", strEmojiMedal, element.IdDiscord)
		content += fmt.Sprintf("\nTotal\t: **%s**", helper.ConvertSecondToFormatHourMinuteSecond(element.Tracked))

		if counter != 3 {
			content += "\n"
		}

		if counter == 3 {
			break
		}
	}

	content += "\n\nPertahankan terus prestasinya dan bagi yang tidak menang pada periode ini jangan berkecil hati karena masih ada periode berikutnya."

	var idChannelDiscord string
	if appEnv == "production" {
		idChannelDiscord = helper.LoadEnvVariable(configs.IdChannelDiscordProduction)
	} else {
		idChannelDiscord = helper.LoadEnvVariable(configs.IdChannelDiscordDevelopment)
	}
	discordRemoteDataSource := datasource.NewDiscordRemoteDataSource(requestDiscord)

	// kirim gif
	isSuccessSendGif := discordRemoteDataSource.SendMessageToChannel(
		idChannelDiscord,
		discord.DiscordCreateMessageBody{
			Content: "https://tenor.com/view/the-rock-entrance-wwe-champion-wwe-raw-gif-12577904",
			Embeds:  []discord.ItemEmbedDiscord{},
		})
	if isSuccessSendGif {
		helper.PrintLog("monthly cattr gif success")
	} else {
		helper.PrintLog("monthly cattr gif failure")
	}

	// kirim info laporannya
	isSuccessSendMessageMonthlyCattr := discordRemoteDataSource.SendMessageToChannel(
		idChannelDiscord,
		discord.DiscordCreateMessageBody{
			Content: content,
			Embeds:  []discord.ItemEmbedDiscord{},
		})
	if isSuccessSendMessageMonthlyCattr {
		helper.PrintLog("monthly cattr success")
	} else {
		helper.PrintLog("monthly cattr failure")
	}
}
