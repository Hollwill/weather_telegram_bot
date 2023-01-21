package settings

import (
	"log"
	"os"

	"github.com/HollWill/weather_telegram_bot/db/models"
	"github.com/HollWill/weather_telegram_bot/db/repositories"
	"github.com/HollWill/weather_telegram_bot/settings"
	"github.com/jmoiron/sqlx"
)

var (
	BotToken        string
	WeatherAPIToken string
	Sdb             *sqlx.DB
	UserRepo        models.UserRepository
)

func initDb() {
	sdb, err := sqlx.Connect("sqlite3", "weather.db")
	if err != nil {
		log.Fatalln(err)
	}
	Sdb = sdb
}

func initBotToken() {
	val, ok := os.LookupEnv("BOT_TOKEN")
	if ok {
		BotToken = val
	} else {
		log.Fatalln("Declare BOT_TOKEN in environment variable")
	}
}

func initWeatherApiToken() {
	val, ok := os.LookupEnv("WEATHER_API_TOKEN")
	if ok {
		WeatherAPIToken = val
	} else {
		log.Fatalln("Declare WEATHER_API_TOKEN in environment variable")
	}
}

func init() {
	initBotToken()
	initDb()
	initWeatherApiToken()

	UserRepo = repositories.NewSqlUserRepository(settings.Sdb)
}
