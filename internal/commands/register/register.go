package register

import (
	"bot_tg/internal/client/telegram"
	user_list "bot_tg/internal/client/user-list"
	"bot_tg/internal/datastruct"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

const (
	registerMessageGood    = "пользователь %s зарегестрирован"
	registerMessageNotGood = "пользователь %s уже зарегестрирован"
)

//RegisterUser implements register command
type RegisterUser interface {
	Start(tg.Update)
}

type registerUser struct {
	chBot    telegram.Client
	chTG     tg.UpdatesChannel
	userList user_list.UserList
}

func RegRegisterUserCommand(chTH tg.UpdatesChannel, chBot telegram.Client, list user_list.UserList) RegisterUser {
	return &registerUser{
		chBot:    chBot,
		chTG:     chTH,
		userList: list,
	}
}

func (r *registerUser) Start(update tg.Update) {
	_, exists := r.userList.GetUser(update.Message.From.ID)
	if !exists {
		r.userList.InsertUser(datastruct.User{
			UserID:   update.Message.From.ID,
			UserName: update.Message.From.UserName,
			Role:     datastruct.Roles[datastruct.Default],
		})
		err := r.chBot.SendMessage(update.Message.Chat.ID,
			fmt.Sprintf(registerMessageGood, update.Message.From.UserName))
		if err != nil {
			log.Println(err)
		}
	} else {
		err := r.chBot.SendMessage(update.Message.Chat.ID,
			fmt.Sprintf(registerMessageNotGood, update.Message.From.UserName))
		if err != nil {
			log.Println(err)
		}
	}
}
