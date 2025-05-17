package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"log"
)

func initConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	return viper.ReadInConfig()
}

func main() {

	if err := initConfig(); err != nil {
		log.Fatal("Ошибка чтения конфига:", err)
	}

	bot, err := tgbotapi.NewBotAPI(viper.GetString("telegram_token"))
	if err != nil {
		log.Fatal("Ошибка создания бота:", err)
	}

	bot.Debug = true
	log.Printf("Бот запущен: @%s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				msg.Text = "Привет! Отправь мне название города, и я пришлю погоду."
			case "weather":
				msg.Text = "Напиши город после команды, например: /weather Москва"
			}
		} else {
			weather, err := GetWeather(update.Message.Text, viper.GetString("weatherstack_api_key"))
			if err != nil {
				msg.Text = "Город не найден. Попробуй еще раз!"
			} else {
				msg.Text = weather
			}
		}
		if err != nil {
			msg.Text = "Город не найден. Попробуйте написать на английском (например, 'Moscow')."
		}

		if _, err := bot.Send(msg); err != nil {
			log.Println("Ошибка отправки сообщения:", err)
		}
	}

}
