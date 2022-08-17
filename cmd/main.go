package main

import (
	"kitty-bot/cmd/app"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/go-resty/resty/v2"
)

func main() {
	s := gocron.NewScheduler(time.Local)
	request := resty.New().R()

	// daily reminder scrum
	s.Every(1).Day().At("16:28").Do(func() {
		app.StartDailyScrum(request)
	})

	// start cron
	// s.StartBlocking()
}
