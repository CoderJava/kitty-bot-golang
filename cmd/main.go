package main

import (
	"kitty-bot/cmd/app"
	"kitty-bot/internal/helper"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/go-resty/resty/v2"
)

func main() {
	s := gocron.NewScheduler(time.Local)
	request := resty.New().R()
	cacheHelper := helper.NewCacheHelper()

	// daily reminder scrum
	s.Every(1).Day().At("16:28").Do(func() {
		app.StartDailyScrum(request)
	})

	// daily hubstaff
	s.Every(1).Day().At("09:30").Do(func() {
		app.StartDailyHubstaff(request, *cacheHelper)
	})

	// start cron
	helper.PrintLog("Running...")
	s.StartBlocking()
}
