package handler

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) handleCommand(update tgbotapi.Update) {
	command := update.Message.Command()

	switch command {
	case "start":
		h.StartCommand(update)
	case "help":
		h.HelpCommand(update)
	case "time":
		h.TimeCommand(update)
	default:
		h.UnknownCommand(update)
	}
}

func (h *Handler) StartCommand(update tgbotapi.Update) {
	msgText := `Hello! This is %s! 
	I'm here to help you find out the weather in any city you wish.
	
	Just type the name of the city, for example: Berlin
	
	If you need help, type /help`

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(msgText, h.bot.Self.FirstName))
	h.bot.Send(msg)
}

func (h *Handler) HelpCommand(update tgbotapi.Update) {
	msgText := `Command Instructions
	
	/start - Start working
	/help - Command instructions
	/time - Current time in chosen city (in process...)`

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	h.bot.Send(msg)
}

func (h *Handler) TimeCommand(update tgbotapi.Update) {

}

func (h *Handler) UnknownCommand(update tgbotapi.Update) {
	msgText := `Sorry! Your command is invalid. Look /help.`
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	h.bot.Send(msg)
}
