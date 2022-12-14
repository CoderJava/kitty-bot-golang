package main

import (
	"kitty-bot/cmd/app"
	"kitty-bot/configs"
	"kitty-bot/internal/helper"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/go-resty/resty/v2"
)

func main() {
	appEnv := helper.LoadEnvVariable(configs.AppEnv)
	s := gocron.NewScheduler(time.Local)
	cacheHelper := helper.NewCacheHelper()
	mh := helper.NewMiddlewareHelper(*cacheHelper)

	httpClientDiscord := resty.New()
	httpClientDiscord.OnBeforeRequest(mh.OnBeforeRequestDiscord)

	httpClientHubstaffAuth := resty.New()
	httpClientHubstaffAuth.OnBeforeRequest(mh.OnBeforeRequestHubstaffAuth)

	httpClientHubstaff := resty.New()
	httpClientHubstaff.OnBeforeRequest(mh.OnBeforeRequestHubstaff)

	httpClientCattr := resty.New()
	httpClientCattr.OnBeforeRequest(mh.OnBeforeRequestCattr)

	if appEnv == "development" {
		httpClientDiscord.SetDebug(true)
		httpClientHubstaffAuth.SetDebug(true)
		httpClientHubstaff.SetDebug(true)
	}

	requestDiscord := httpClientDiscord.R()
	requestCattr := httpClientCattr.R()

	// daily reminder scrum at weekday
	s.Every(1).Day().At("16:28").Do(func() {
		app.StartDailyScrum(
			requestDiscord,
			[]string{"Mon", "Tue", "Wed", "Thu", "Fri"},
		)
	})

	// daily reminder scrum at saturday
	s.Every(1).Week().Saturday().At("11:28").Do(func() {
		app.StartDailyScrum(
			requestDiscord,
			[]string{"Sat"},
		)
	})

	// daily cattr at 09:30 on every day
	s.Every(1).Day().At("09:30").Do(func() {
		app.StartDailyCattr(
			requestCattr,
			requestDiscord,
			*cacheHelper,
		)
	})

	// monthly cattr at 10:30 on day of month 27
	s.Cron("30 10 27 * *").Do(func() {
		app.StartMonthlyCattr(
			requestCattr,
			requestDiscord,
			*cacheHelper,
		)
	})

	// start cron
	helper.PrintLog("Running...")
	s.StartBlocking()
}
