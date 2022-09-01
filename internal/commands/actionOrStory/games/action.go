package games

import (
	"bot_tg/internal/client/telegram"
	user_list "bot_tg/internal/client/user-list"
	"bot_tg/internal/datastruct"
	"bot_tg/internal/utils"
	"context"
	"fmt"
	"github.com/enescakir/emoji"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strconv"
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
	rules = "У игроков есть команда \n```/предложить {текст предложения}```\nПри новом предложении - старое перетирается\nЧерез 10 минут бот запустит голосование, где необходимо будет выбрать вариант\n"
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

func (r *pidGameAction) Start(
	ctx context.Context,
	wg *sync.WaitGroup,
	update tg.Update,
	updateCh tg.UpdatesChannel) {

	value := ctx.Value(CurrPlayer)
	if player, ok := value.(datastruct.User); ok {

		err := r.chBot.SendMessage(update.Message.Chat.ID, rules)
		if err != nil {
			log.Println(err)
		}

		offer := make(map[int]*datastruct.Vote)

		for {
			select {
			case currUpdate := <-updateCh:
				msgSplit := strings.Split(currUpdate.Message.Text, " ")
				msgText := strings.Join(msgSplit[1:], " ")
				switch msgSplit[0] {
				case "/предложить":
					if r.userList.Exists(currUpdate.Message.From.ID) &&
						currUpdate.Message.From.ID != player.UserID {
						offer[currUpdate.Message.From.ID] = &datastruct.Vote{
							UserID:    0,
							OfferText: msgText,
							IsVoted:   false,
						}
						err = r.chBot.SendMessage(currUpdate.Message.Chat.ID, "Предложение принято")
						if err != nil {
							log.Println(err)
						}
					}
				}
			case <-time.After(time.Minute * 1):
				r.startStepOne(ctx, wg, updateCh, offer, update.Message.Chat.ID)
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
	votes map[int]*datastruct.Vote,
	chatID int64) {

	var (
		votesOffers  string
		offersNumber = 1
	)

	offersFinal := make([]datastruct.VoteOffer, 0, len(votesOffers))
	keyboardButtons := make([]string, 0)

	for _, v := range votes {
		votesOffers += fmt.Sprintf("Предложение #%d:\n%s\n", offersNumber, v.OfferText)
		offersFinal = append(offersFinal, datastruct.VoteOffer{
			OfferText:   v.OfferText,
			VotesCount:  0,
			OfferNumber: offersNumber,
		})
		keyboardButtons = append(keyboardButtons, strconv.Itoa(offersNumber))
		offersNumber++
	}

	keyboard := utils.CreateCustomKeyboardByThreeButtons(keyboardButtons)

	err := r.chBot.SendMessage(chatID, "Начинается этап голосования, список предложений:")
	if err != nil {
		log.Println(err)
	}

	err = r.chBot.SendMessage(chatID, votesOffers)
	if err != nil {
		log.Println(err)
	}

	err = r.chBot.SendMessageWithKeyboard(chatID, "Пункты голосования:", *keyboard)
	if err != nil {
		log.Println(err)
	}

	for {
		select {
		case upd := <-updateCh:
			if vote, ok := votes[upd.Message.From.ID]; ok {
				voteNumber, errVote := strconv.Atoi(upd.Message.Text)
				if errVote == nil && voteNumber >= 1 && voteNumber <= offersNumber && vote.IsVoted == false {
					votes[upd.Message.From.ID].IsVoted = true
					for _, off := range offersFinal {
						if off.OfferNumber == voteNumber {
							off.VotesCount++
						}
					}
				}
			}
		case <-time.After(time.Second * 10):
			r.startStepTwo(ctx, wg, updateCh, offersFinal, chatID)
		}
	}
}

func (r *pidGameAction) startStepTwo(
	ctx context.Context,
	wg *sync.WaitGroup,
	updateCh tg.UpdatesChannel,
	result []datastruct.VoteOffer,
	chatID int64) {
	resultMessage := "Результаты голосования:"
	for _, v := range result {
		voteOffer := fmt.Sprintf("%d: ", v.OfferNumber)
		for i := 1; i < v.VotesCount; i++ {
			voteOffer += fmt.Sprintf("%s ", emoji.ThumbsUp)
		}
		voteOffer += "\n"
		resultMessage += voteOffer
	}

	err := r.chBot.SendMessage(chatID, resultMessage)
	if err != nil {
		log.Println(err)
	}
}
