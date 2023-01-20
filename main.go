package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"

	"github.com/HollWill/weather_telegram_bot/db/sqlstore"
	"github.com/HollWill/weather_telegram_bot/handlers"
	"github.com/HollWill/weather_telegram_bot/mailing"
	"github.com/HollWill/weather_telegram_bot/predicates"
)

var botToken string

func initLog() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
}

func init() {
	val, ok := os.LookupEnv("BOT_TOKEN")
	if ok {
		botToken = val
	} else {
		log.Fatalln("Declare BOT_TOKEN in environment variable")
	}
	sdb, err := sqlx.Connect("sqlite3", "weather.db")
	if err != nil {
		log.Fatalln(err)
	}
	sqlstore.CreateTables(sdb)
}

func main() {
	bot, err := telego.NewBot(botToken, telego.WithDefaultLogger(false, false))
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
