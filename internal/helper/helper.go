package helper

import (
	"fmt"
	"kitty-bot/cmd/domain/hubstaff"
	"kitty-bot/configs"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Untuk mengecek apakah argumen search terdapat didalam argumen values
func ContainString(values []string, search string) bool {
	for _, value := range values {
		if value == search {
			return true
		}
	}
	return false
}

// Untuk mencetak nilai message ke console dengan tambahan info "2006-01-02 15:04:05"
func PrintLog(message string) {
	now := time.Now()
	formattedNow := now.Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] %s\n", formattedNow, message)
}

// Untuk membuat string dengan info "2006-01-02 15:04:05" ditambah dengan argumen message
func SprintLog(message string) string {
	now := time.Now()
	formattedNow := now.Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s] %s\n", formattedNow, message)
}

// Untuk mengambil nilai dari environment variable sesuai dengan key yang diberikan
func LoadEnvVariable(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv(key)
}

// Untuk konversi detik menjadi HH:mm:ss
func ConvertSecondToFormatHourMinuteSecond(second int) string {
	var hour int
	var minute int
	if second%3600 > 0 {
		hour = second / 3600
		second -= hour * 3600
	}
	if second/60 > 0 {
		minute = second / 60
		second -= minute * 60
	}
	strHour := fmt.Sprintf("%02d", hour)
	strMinute := fmt.Sprintf("%02d", minute)
	strSecond := fmt.Sprintf("%02d", second)
	return fmt.Sprintf("%s:%s:%s", strHour, strMinute, strSecond)
}

func GetNameAndIdDiscordByIdHubstaff(idHubstaff string) (name string, idDiscord string) {
	var keyIdDiscord string
	switch idHubstaff {
	case LoadEnvVariable(configs.IdHubstaffYudiSetiawan):
		keyIdDiscord = configs.IdDiscordYudiSetiawan
		name = configs.NameYudiSetiawan
	case LoadEnvVariable(configs.IdHubstaffRyan):
		keyIdDiscord = configs.IdDiscordRyan
		name = configs.NameRyanAlfarisi
	case LoadEnvVariable(configs.IdHubstaffSabrino):
		keyIdDiscord = configs.IdDiscordSabrino
		name = configs.NameSabrino
	case LoadEnvVariable(configs.IdHubstaffRioDwi):
		keyIdDiscord = configs.IdDiscordRioDwi
		name = configs.NameRioDwiPrabowo
	case LoadEnvVariable(configs.IdHubstaffBobby):
		keyIdDiscord = configs.IdDiscordBobby
		name = configs.NameAdhityaBobby
	case LoadEnvVariable(configs.IdHubstaffAditama):
		keyIdDiscord = configs.IdDiscordAditama
		name = configs.NameAditama
	case LoadEnvVariable(configs.IdHubstaffAldoFaiz):
		keyIdDiscord = configs.IdDiscordAldoFaizi
		name = configs.NameAldoFaizi
	case LoadEnvVariable(configs.IdHubstaffDewi):
		keyIdDiscord = configs.IdDiscordDewi
		name = configs.NameDewiLilian
	case LoadEnvVariable(configs.IdHubstaffAbdulAziz):
		keyIdDiscord = configs.IdDiscordAbdulAziz
		name = configs.NameAbdulAziz
	case LoadEnvVariable(configs.IdHubstaffRianto):
		keyIdDiscord = configs.IdDiscordRianto
		name = configs.NameRianto
	}
	if keyIdDiscord != "" {
		idDiscord = LoadEnvVariable(keyIdDiscord)
	}
	return
}

// Untuk mem-filter slice x yang terpenuhi berdasarkan kriteria parameter isFiltered
func FilterTemplateMessageHubstaff(
	x []hubstaff.TemplateMessageHubstaff,
	isFiltered func(int) bool,
) (resultFilter []hubstaff.TemplateMessageHubstaff) {
	for index, element := range x {
		if isFiltered(index) {
			resultFilter = append(resultFilter, element)
		}
	}
	return
}

// Untuk mencari posisi index dari sebuah slice jika parameter predicate terpenuhi.
// Jika tidak ada yang terpenuhi maka, nilai kembaliannya adalah -1.
func SliceIndex(limit int, predicate func(int) bool) int {
	for index := 0; index < limit; index++ {
		if predicate(index) {
			return index
		}
	}
	return -1
}
