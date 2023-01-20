package mailing

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/HollWill/weather/api_services"
	"github.com/HollWill/weather/structures"
	"github.com/HollWill/weather_telegram_bot/db/models"
	"github.com/HollWill/weather_telegram_bot/db/repositories"
	"github.com/go-co-op/gocron"
	"github.com/jmoiron/sqlx"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

var (
	userRepo        models.UserRepository
	weatherAPIToken string
	Sheduler        *gocron.Scheduler
)

func init() {
	val, ok := os.LookupEnv("WEATHER_API_TOKEN")
	if ok {
		weatherAPIToken = val
	} else {
		log.Fatalln("Declare WEATHER_API_TOKEN in environment variable")
	}
	sdb, err := sqlx.Connect("sqlite3", "weather.db")
	if err != nil {
		log.Fatalln(err)
	}
	userRepo = repositories.NewSqlUserRepository(sdb)
}

func sendMessage(bot *telego.Bot, user models.User) func() {
	return func() {
		service := api_services.WeatherApiComService{
			Coords: structures.Coords{
				Lat:  float64(user.Latitude),
				Long: float64(user.Longitude),
			},
			ApiKey: weatherAPIToken,
		}
		parser := service.GetParser()

		weather := parser.Parse(service.FetchData())

		_, _ = bot.SendMessage(tu.Message(
			tu.ID(int64(user.ID)),
			fmt.Sprint(weather),
		))
	}
}

func Mailing(bot *telego.Bot) {
	if Sheduler != nil {
		Sheduler.Clear()
	}
	users, err := userRepo.GetAll(context.Background())
	if err != nil {
		log.Println(err)
	}

	Sheduler = gocron.NewScheduler(time.UTC)

	for _, user := range users {
		Sheduler.Cron(user.CronTab).Do(sendMessage(bot, user))
	}

	Sheduler.StartAsync()
}
