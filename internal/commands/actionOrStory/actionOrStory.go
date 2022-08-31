package actionOrStory

import (
	"bot_tg/internal/client/telegram"
	user_list "bot_tg/internal/client/user-list"
	"bot_tg/internal/commands/actionOrStory/games"
	"context"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"sync"
)

const (
	truth             = "правда"
	action            = "действие"
	msgBegin          = "Игра началась...\n@%s Выбирает 'правда' или 'действие'"
	msgNotEnoughUsers = "недостаточно участников для старта игры"
)

//CurrPlayer context player key
const CurrPlayer = 0

var actionKeyBoard = tg.NewReplyKeyboard(
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton(truth),
		tg.NewKeyboardButton(action),
	),
)

type PIDGame interface {
	Start(*sync.WaitGroup, tg.Update, tg.UpdatesChannel)
}

type pidGame struct {
	chBot    telegram.Client
	chTG     tg.UpdatesChannel
	userList user_list.UserList
}

func RegNewGameCommand(chTH tg.UpdatesChannel, chBot telegram.Client, list user_list.UserList) PIDGame {
	return &pidGame{
		chBot:    chBot,
		chTG:     chTH,
		userList: list,
	}
}

func (r *pidGame) Start(wg *sync.WaitGroup, update tg.Update, updateCh tg.UpdatesChannel) {
	player, exists := r.userList.GetRandomUser()
	if !exists {
		err := r.chBot.SendMessage(update.Message.Chat.ID, msgNotEnoughUsers)
		if err != nil {
			log.Println(err)
		}
		wg.Done()
		return
	}

	ctx := context.WithValue(context.Background(), CurrPlayer, player)

	err := r.chBot.SetAvatar(update.Message.Chat.ID, telegram.GameStartPic)
	if err != nil {
		log.Println(err)
	}

	str := fmt.Sprintf(msgBegin, player.UserName)

	err = r.chBot.SendMessageWithKeyboard(update.Message.Chat.ID, str, actionKeyBoard)
	if err != nil {
		log.Println(err)
	}

	for {
		select {
		case currUpdate := <-updateCh:
			if currUpdate.Message.From.ID == player.UserID {
				switch currUpdate.Message.Text {
				case action:
					aGame := games.RegNewGameAction(r.chTG, r.chBot, r.userList)
					aGame.Start(ctx, wg, currUpdate, updateCh)
				case truth:

				}
			} else {
				if currUpdate.Message.Text == action || currUpdate.Message.Text == truth {
					err = r.chBot.SendMessage(currUpdate.Message.Chat.ID, "только избранный может выбирать")
					if err != nil {
						log.Println(err)
					}
				}
			}
		default:
		}
	}
}
