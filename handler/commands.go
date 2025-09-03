package handler

import (
	"fmt"
	"log"
	"math"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) handleCommand(update tgbotapi.Update) {
	command := update.Message.Command()

	switch command {
	case "start":
		h.handleStartCommand(update)
	case "help":
		h.handleHelpCommand(update)
	case "current":
		h.handleCurrentCommand(update)
	default:
		h.handleUnknownCommand(update)
	}
}

func (h *Handler) handleStartCommand(update tgbotapi.Update) {
	msgText := `Hello! This is %s! 
	I'm here to help you find out the weather in any city you wish.
	
	Just type the name of the city, for example: Berlin
	
	If you need help, type /help`

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(msgText, h.bot.Self.FirstName))
	h.bot.Send(msg)
}

func (h *Handler) handleHelpCommand(update tgbotapi.Update) {
	msgText := `Command Instructions
	
	/start - Start working
	/help - Command instructions
	/current - Weather in current place (in process...)`

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	h.bot.Send(msg)
}

func (h *Handler) handleUnknownCommand(update tgbotapi.Update) {
	msgText := `Sorry! Your command is invalid. Look /help.`
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	h.bot.Send(msg)
}

func (h *Handler) handleCurrentCommand(update tgbotapi.Update) {
	initMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Send me your geolocation, and I will tell you what the weather is like at your point!")
	h.bot.Send(initMsg)

	userID := update.Message.From.ID
	h.usersState[userID] = "waiting_location"

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "📍 Please, send me your location!")
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButtonLocation("Send my location")))
	h.bot.Send(msg)
}

func (h *Handler) handleLocation(update tgbotapi.Update) {
	userID := update.Message.From.ID
	defer delete(h.usersState, userID)

	if update.Message.Location == nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			"Please, send me location, nothing else!\n"+
				"Use \"Send my location\" button!\n"+
				"If you want to try again, click /current")
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		h.bot.Send(msg)
		return
	}

	loc := update.Message.Location

	weather, err := h.owClient.Weather(loc.Latitude, loc.Longitude)
	if err != nil {
		log.Println(err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Couldn't get weather!")
		h.bot.Send(msg)
		return
	}

	resp := fmt.Sprintf("Weather in your location: %d", int(math.Round(weather.Temp)))
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, resp)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

	h.bot.Send(msg)
}
