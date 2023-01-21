package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"

	"github.com/HollWill/weather_telegram_bot/db/sqlstore"
	"github.com/HollWill/weather_telegram_bot/handlers"
	"github.com/HollWill/weather_telegram_bot/mailing"
	"github.com/HollWill/weather_telegram_bot/predicates"
	"github.com/HollWill/weather_telegram_bot/settings"
)

func initLog() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
}

func init() {
	sqlstore.CreateTables(settings.Sdb)
	initLog()
}

func main() {
	bot, err := telego.NewBot(settings.BotToken, telego.WithDefaultLogger(false, false))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	go mailing.Mailing(bot)

	updates, _ := bot.UpdatesViaLongPulling(nil)
	defer bot.StopLongPulling()

	bh, _ := th.NewBotHandler(bot, updates)

	bh.Handle(handlers.StartHandler, th.CommandEqual("start"))

	bh.Handle(handlers.WeatherHandler, th.CommandEqual("weather"))

	bh.Handle(handlers.LocationHandler, predicates.HasLocation())

	bh.Handle(handlers.CrontabHandler, th.CommandEqualArgc("crontab", 5))

	defer bh.Stop()

	bh.Start()
}
