package games

import (
	"bot_tg/internal/client/telegram"
	user_list "bot_tg/internal/client/user-list"
	"bot_tg/internal/datastruct"
	"context"
	"fmt"
	"github.com/enescakir/emoji"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
	"sync"
	"time"
)

//context keys
const (
	CurrPlayer = iota
	CurrChatID
)

const (
	wishAction   = "/загадать действие"
	startMessage = "%s выбрал действие, рейтинг игрока: %f\nСписок игроков, принимающих участие в предложении действия: %s"
	rules        = "У игроков из списка предлагающих есть команда \n```/предложить {текст предложения}```\nПри новом предложении - старое перетирается\nЧерез 10 минут бот запустит голосование, где необходимо будет выбрать вариант\n"
)

type PIDAction interface {
	Start(context.Context, *sync.WaitGroup, tg.Update, tg.UpdatesChannel)
}

type pidGameAction struct {
	chBot    telegram.Client
	chTG     tg.UpdatesChannel
	userList user_list.UserList
}

func RegNewGameAction(chTH tg.UpdatesChannel, chBot telegram.Client, list user_list.UserList) PIDAction {
	return &pidGameAction{
		chBot:    chBot,
		chTG:     chTH,
		userList: list,
	}
}

func (r *pidGameAction) Start(ctx context.Context, wg *sync.WaitGroup, update tg.Update, updateCh tg.UpdatesChannel) {
	value := ctx.Value(CurrPlayer)
	if player, ok := value.(datastruct.User); ok {
		votingUsers := r.userList.GetAllUsers()
		for i, user := range votingUsers {
			if user.UserID == player.UserID {
				votingUsers[i] = votingUsers[len(votingUsers)-1]
				votingUsers = votingUsers[:len(votingUsers)-1]
			}
		}

		voting := make(map[int]datastruct.Vote, len(votingUsers))
		var strVotingUsers string
		for _, user := range votingUsers {
			strVotingUsers += fmt.Sprintf("%s %s\n", user.UserName, emoji.Mushroom)
			voting[user.UserID] = datastruct.Vote{
				UserID:    user.UserID,
				MessageID: 0,
				IsVoted:   false,
			}
		}

		err := r.chBot.SendMessage(update.Message.Chat.ID, fmt.Sprintf(startMessage, player.UserName, player.Rating, strVotingUsers))
		if err != nil {
			log.Println(err)
		}
		err = r.chBot.SendMessage(update.Message.Chat.ID, rules)
		if err != nil {
			log.Println(err)
		}

		for {
			select {
			case currUpdate := <-updateCh:
				msgSplit := strings.Split(currUpdate.Message.Text, " ")
				msgText := strings.Join(msgSplit[1:], " ")
				switch msgSplit[0] {
				case "/предложить":
					if vote, ex := voting[update.Message.From.ID]; ex {
						vote.OfferTxt = msgText
						voting[update.Message.From.ID] = vote
						err = r.chBot.SendMessage(update.Message.Chat.ID, fmt.Sprintf("Предложение ```%s``` принято", msgText))
						if err != nil {
							log.Println(err)
						}
					} else {
						err = r.chBot.SendMessage(update.Message.Chat.ID, "Предлагать не могут незарегестрированные/избранный")
						if err != nil {
							log.Println(err)
						}
					}
				}
			case <-time.After(time.Minute * 10):
				ctxNext := context.WithValue(ctx, CurrChatID, update.Message.Chat.ID)
				r.startStepOne(ctxNext, wg, updateCh, voting)
			}
		}
	} else {
		log.Println("can't catch user from ctx")
		return
	}
}

func (r *pidGameAction) startStepOne(
	ctx context.Context,
	wg *sync.WaitGroup,
	updateCh tg.UpdatesChannel,
	votes map[int]datastruct.Vote) {
	var votesOffers string
	var i int
	for _, v := range votes {
		i++
		votesOffers += fmt.Sprintf("Предложение #%d:\n%s\n", i, v.OfferTxt)
	}
	var chatID int64
	value := ctx.Value(CurrChatID)
	if chID, ok := value.(int64); ok {
		chatID = chID
	}

	err := r.chBot.SendMessage(chatID, "Начинается этап голосования, список предложений:")
	if err != nil {
		log.Println(err)
	}

	_, err = r.chBot.SendMessageAndReturnMSG(chatID, votesOffers)
	if err != nil {
		log.Println(err)
	}

	err = r.chBot.SendMessage(chatID, "Чтобы проголосовать выберите на клавиаутер номер сообщения.")
	if err != nil {
		log.Println(err)
	}
}
