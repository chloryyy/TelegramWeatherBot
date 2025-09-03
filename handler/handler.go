package handler

import (
	"fmt"
	"log"
	"math"

	"github.com/chloryyy/WeatherBot/clients/openweather"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	bot        *tgbotapi.BotAPI
	owClient   *openweather.OpenWeatherClient
	usersState map[int64]string
}

func New(bot *tgbotapi.BotAPI, owClient *openweather.OpenWeatherClient) *Handler {
	return &Handler{
		bot:        bot,
		owClient:   owClient,
		usersState: make(map[int64]string),
	}
}

func (h *Handler) handleUpdate(update tgbotapi.Update) {

	if update.Message == nil {
		return
	}

	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	userID := update.Message.From.ID

	if state, ok := h.usersState[userID]; ok {
		switch state {
		case "waiting_location":
			h.handleLocation(update)
			return
		}
	}
	if update.Message.IsCommand() {
		h.handleCommand(update)
		return
	}

	// If we got a message

	coordinates, err := h.owClient.Coordinates(update.Message.Text)
	if err != nil {
		log.Println(err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не смогли обработать сообщение")
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}

	weather, err := h.owClient.Weather(coordinates.Lat, coordinates.Lon)
	if err != nil {
		log.Println(err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не смогли узнать погоду в этой местности")
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		fmt.Sprintf("Температура в местности %s: %d°C", update.Message.Text, int(math.Round(weather.Temp))),
	)
	msg.ReplyToMessageID = update.Message.MessageID

	h.bot.Send(msg)

}

func (h *Handler) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := h.bot.GetUpdatesChan(u)

	for update := range updates {
		h.handleUpdate(update)
	}
}
