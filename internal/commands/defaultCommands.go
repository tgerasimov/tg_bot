package commands

import (
	"bot_tg/internal/client/telegram"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

type defaultCommands struct {
	chBot telegram.Client
	chTG  tg.UpdatesChannel
}

func RegDefaultCommands(chTG tg.UpdatesChannel, chBot telegram.Client) *defaultCommands {
	return &defaultCommands{chBot: chBot, chTG: chTG}
}

func (d *defaultCommands) CommandList(update tg.Update) {
	for k, v := range CommandList {
		err := d.chBot.SendMessage(update.Message.Chat.ID,
			fmt.Sprintf("команда: '%s'\nфункционал: '%s'", k, v))
		if err != nil {
			log.Println(err)
		}
	}
}
