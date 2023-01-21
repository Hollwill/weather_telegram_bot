package handlers

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/HollWill/weather/api_services"
	"github.com/HollWill/weather/structures"
	"github.com/HollWill/weather_telegram_bot/db/models"
	"github.com/HollWill/weather_telegram_bot/mailing"
	"github.com/HollWill/weather_telegram_bot/settings"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func StartHandler(bot *telego.Bot, update telego.Update) {
	keyboard := tu.Keyboard(
		tu.KeyboardRow(
			tu.KeyboardButton("Отправить геолокацию").WithRequestLocation(),
		),
	)

	_, _ = bot.SendMessage(tu.Messagef(
		tu.ID(update.Message.Chat.ID),
		"Привет, %s! \n Отправь свою локацию чтобы я знал где ты находишься",
		update.Message.From.FirstName,
	).WithReplyMarkup(keyboard),
	)

	u := models.User{}
	u.ID = int(update.Message.Chat.ID)
	u.Name = update.Message.From.FirstName
	ctx := context.Background()
	err := settings.UserRepo.Save(ctx, &u)
	if err != nil {
		log.Println(err)
	}
}

func LocationHandler(bot *telego.Bot, update telego.Update) {
	u, err := settings.UserRepo.FindByID(context.Background(), int(update.Message.Chat.ID))
	if err != nil {
		log.Println("User not found")
		_, _ = bot.SendMessage(
			tu.Message(
				tu.ID(update.Message.Chat.ID),
				"Пользователь не найден. Начните с команды /start",
			),
		)

	}

	u.Latitude = float32(update.Message.Location.Latitude)
	u.Longitude = float32(update.Message.Location.Longitude)
	settings.UserRepo.Save(context.Background(), u)

	_, _ = bot.SendMessage(
		tu.Message(
			tu.ID(update.Message.Chat.ID),
			`Геолокация успешно сохранена.
			Введите время для получения рассылки в формате /crontab * * * * 
			Например '/crontab 1 * * * *' для получения рассылки каждый час в 1 минуту.`,
		).WithReplyMarkup(&telego.ReplyKeyboardRemove{RemoveKeyboard: true}),
	)
}

func CrontabHandler(bot *telego.Bot, update telego.Update) {
	_, args := tu.ParseCommand(update.Message.Text)
	strings.Join(args, " ")

	u, err := settings.UserRepo.FindByID(context.Background(), int(update.Message.Chat.ID))
	if err != nil {
		log.Println(err)
	}

	u.CronTab = strings.Join(args, " ")
	settings.UserRepo.Save(context.Background(), u)

	go mailing.Mailing(bot)
	_, _ = bot.SendMessage(
		tu.Message(
			tu.ID(update.Message.Chat.ID),
			`Время рассылки успешно сохранено.`,
		).WithReplyMarkup(&telego.ReplyKeyboardRemove{RemoveKeyboard: true}),
	)
}

func WeatherHandler(bot *telego.Bot, update telego.Update) {
	u, err := settings.UserRepo.FindByID(context.Background(), int(update.Message.Chat.ID))
	if err != nil {
		log.Println("Пользователь не найден.")
		_, _ = bot.SendMessage(
			tu.Message(
				tu.ID(update.Message.Chat.ID),
				"Пользователь не найден. Начните с команды /start",
			),
		)
	} else {
		service := api_services.WeatherApiComService{
			Coords: structures.Coords{Lat: float64(u.Latitude), Long: float64(u.Longitude)},
			ApiKey: settings.WeatherAPIToken,
		}
		parser := service.GetParser()

		weather := parser.Parse(service.FetchData())

		_, _ = bot.SendMessage(tu.Message(
			tu.ID(update.Message.Chat.ID),
			fmt.Sprint(weather),
		))

	}
}
