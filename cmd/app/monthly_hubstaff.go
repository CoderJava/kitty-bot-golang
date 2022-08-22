package app

import (
	"fmt"
	"kitty-bot/cmd/datasource"
	"kitty-bot/cmd/domain/discord"
	"kitty-bot/cmd/domain/hubstaff"
	"kitty-bot/configs"
	"kitty-bot/internal/helper"
	"sort"
	"time"

	"github.com/go-resty/resty/v2"
)

func StartMonthlyHubstaff(
	requestHubstaff *resty.Request,
	requestHubstaffAuth *resty.Request,
	requestDiscord *resty.Request,
	cacheHelper helper.CacheHelper,
) {
	// baca APP_ENV
	appEnv := helper.LoadEnvVariable(configs.AppEnv)

	// hit ke endpoint login hubstaff dan simpan access tokennya didalam cache
	hubstaffRemoteDataSource := datasource.NewHubstaffRemoteDataSource(
		requestHubstaff,
		requestHubstaffAuth,
	)
	loginResponse := hubstaffRemoteDataSource.Login()
	if (hubstaff.LoginResponse{}) == loginResponse {
		helper.PrintLog("endpoint login hubstaff gagal")
		return
	}
	cacheHelper.Set(configs.AccessToken, loginResponse.AccessToken)
	cacheHelper.Set(configs.RefreshToken, loginResponse.RefreshToken)

	// ambil waktu saat ini
	// dan tentukan tanggal mulai dan akhir
	now := time.Now()
	var startDate time.Time
	var stopDate time.Time
	getPeriod(now, &startDate, &stopDate)

	// formatting variable startDate dan stopDate
	patternDate := "2006-01-02"
	strStartDate := startDate.Format(patternDate)
	strStopDate := stopDate.Format(patternDate)

	dailyActivities := []hubstaff.ItemDailyActivityResponse{}
	var nextPageStartId int
	for {
		dailyActivityResponse := hubstaffRemoteDataSource.GetDailyActivityByRangeDate(
			strStartDate,
			strStopDate,
			nextPageStartId,
		)
		dailyActivities = append(dailyActivities, dailyActivityResponse.DailyActivities...)
		nextPageStartId = dailyActivityResponse.Pagination.NextPageStartId
		if nextPageStartId == 0 {
			break
		}
	}

	listMessages := []hubstaff.TemplateMessageHubstaff{}
	userIdHubstaffAdeIskandar := helper.LoadEnvVariable(configs.IdHubstaffAdeIskandar)
	for _, itemDailyActivityResponse := range dailyActivities {
		userIdHubstaff := fmt.Sprint(itemDailyActivityResponse.UserId)
		if userIdHubstaff == userIdHubstaffAdeIskandar {
			continue
		}

		name, idDiscord := helper.GetNameAndIdDiscordByIdHubstaff(userIdHubstaff)

		// abaikan datanya jika id discord user-nya tidak valid
		if idDiscord == "" {
			continue
		}

		tracked := itemDailyActivityResponse.InputTrackedInSeconds
		idle := itemDailyActivityResponse.IdleInSeconds

		// cek apakah data user tersebut sudah ada sebelumnya.
		// kalau ada maka, hitung juga
		listAlready := helper.FilterTemplateMessageHubstaff(listMessages, func(index int) bool {
			element := listMessages[index]
			return element.IdHubstaff == userIdHubstaff
		})
		if len(listAlready) > 0 {
			// hitung nilai tracked dan idle terhadap data sebelumnya
			for _, itemAlready := range listAlready {
				tracked += itemAlready.Tracked
				idle += itemAlready.Idle
			}

			// cari index pada data sebelumnya yang didalam list template messages
			indexAlready := helper.SliceIndex(len(listMessages), func(index int) bool {
				element := listMessages[index]
				return element.IdHubstaff == userIdHubstaff
			})
			listMessages[indexAlready] = hubstaff.TemplateMessageHubstaff{
				IdHubstaff: userIdHubstaff,
				IdDiscord:  idDiscord,
				Name:       name,
				Tracked:    tracked,
				Idle:       idle,
			}
		} else {
			listMessages = append(listMessages, hubstaff.TemplateMessageHubstaff{
				IdHubstaff: userIdHubstaff,
				IdDiscord:  idDiscord,
				Name:       name,
				Tracked:    tracked,
				Idle:       idle,
			})
		}
	}

	listBestTracked := []hubstaff.TemplateMessageHubstaff{}
	listBestTracked = append(listBestTracked, listMessages...)

	listBestIdle := []hubstaff.TemplateMessageHubstaff{}
	listBestIdle = append(listBestIdle, listMessages...)

	// urutkan user yang best tracked secara ascending
	sort.Slice(listBestTracked, func(x, y int) bool {
		totalX := listBestTracked[x].Tracked - listBestTracked[y].Idle
		totalY := listBestTracked[y].Tracked - listBestTracked[y].Idle
		return totalX < totalY
	})

	// urutkan user yang best idle secara ascending
	sort.Slice(listBestIdle, func(x, y int) bool {
		return listBestIdle[x].Idle < listBestIdle[y].Idle
	})

	strNowYear := startDate.Format("2006")
	patternDatePeriod := "02 Jan"
	strPeriod := startDate.Format(patternDatePeriod) + " - " + stopDate.Format(patternDatePeriod)
	content := "@everyone Welcome ladies and gentlemen to the"
	content += "\n:trophy: HubStaff Championship " + strNowYear + " :trophy:"
	content += "\n\n"
	content += "Setelah sebulan lamanya para peserta kita bekerja dengan **serius tapi santai** untuk memenangkan kompetisi ini "
	content += "maka, tibalah waktunya untuk mengumumkan the next winner of HubStaff Championship " + strNowYear + " periode **" + strPeriod + "**."

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
		content += fmt.Sprintf("\nTracked\t: **%s**", helper.ConvertSecondToFormatHourMinuteSecond(element.Tracked-element.Idle))

		if counter != 3 {
			content += "\n"
		}

		if counter == 3 {
			break
		}
	}

	// buat pesan untuk idle terbaik
	if listBestIdle[len(listBestIdle)-1].Idle > 0 {
		content += "\n\n"
		content += ":sleeping: **Idle Terbaik**"
		content += "\n============="
		counter = 0
		for index := len(listBestIdle) - 1; index >= 0; index-- {
			counter += 1
			element := listBestIdle[index]
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
			content += fmt.Sprintf("\nIdle\t: **%s**", helper.ConvertSecondToFormatHourMinuteSecond(element.Idle))

			if counter != 3 {
				content += "\n"
			}

			if counter == 3 {
				break
			}
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
		helper.PrintLog("monthly hubstaff gif success")
	} else {
		helper.PrintLog("monthly hubstaff gif failure")
	}

	// kirim info laporannya
	isSuccessSendMessageMonthlyHubstaff := discordRemoteDataSource.SendMessageToChannel(
		idChannelDiscord,
		discord.DiscordCreateMessageBody{
			Content: content,
			Embeds:  []discord.ItemEmbedDiscord{},
		})
	if isSuccessSendMessageMonthlyHubstaff {
		helper.PrintLog("monthly hubstaff success")
	} else {
		helper.PrintLog("monthly hubstaff failure")
	}
}

// untuk menentukan periode hubstaff selama sebulan.
// Intinya dari tgl 26-25.
func getPeriod(now time.Time, startDate *time.Time, stopDate *time.Time) {
	if now.Month() == 1 {
		// ambil ke tahun sebelumnya
		// contoh: 26 Des 2022 - 25 Jan 2023
		*startDate = time.Date(
			now.Year()-1,
			12,
			26,
			0,
			0,
			0,
			0,
			now.Location(),
		)
		*stopDate = time.Date(
			now.Year(),
			now.Month(),
			25,
			0,
			0,
			0,
			0,
			now.Location(),
		)
	} else {
		// ambil pada tahun yang sama
		// contoh: 26 Jan 2022 - 25 Feb 2022
		*startDate = time.Date(
			now.Year(),
			now.Month()-1,
			26,
			0,
			0,
			0,
			0,
			now.Location(),
		)
		*stopDate = time.Date(
			now.Year(),
			now.Month(),
			25,
			0,
			0,
			0,
			0,
			now.Location(),
		)
	}
}
