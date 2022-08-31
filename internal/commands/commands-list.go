package commands

import (
	"bot_tg/internal/client/telegram"
	"bot_tg/internal/commands/actionOrStory"
	"context"
	"sync"

	user_list "bot_tg/internal/client/user-list"
	"bot_tg/internal/commands/register"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	once sync.Once
)

var CommandList = map[string]string{
	"/зарегестрироваться": "команда используется для регистрации пользователя в игре",
	"/ПиД":                "начать игру",
	"/команды":            "показывает список команд",
	"/участники":          "показать всех участников",
}

func RegisterAllCommands(
	ctx context.Context,
	chTG tg.UpdatesChannel,
	chBot telegram.Client,
	userList user_list.UserList,
) {
	once.Do(func() {
		defComma := RegDefaultCommands(chTG, chBot)

		regComma := register.RegRegisterUserCommand(chTG, chBot, userList)

		game := actionOrStory.RegNewGameCommand(chTG, chBot, userList)

		startListenCommands(ctx, chTG, regComma, game, defComma)
	})
}

func startListenCommands(ctx context.Context,
	chTG tg.UpdatesChannel,
	registerUserComma register.RegisterUser,
	game actionOrStory.PIDGame,
	defaultComma *defaultCommands,
) {
	for {
		select {
		case update := <-chTG:
			if _, ok := CommandList[update.Message.Text]; ok {
				switch update.Message.Text {
				case "/зарегестрироваться":
					registerUserComma.Start(update)
				case "/команды":
					defaultComma.CommandList(update)
				case "/ПиД":
					wg := new(sync.WaitGroup)
					wg.Add(1)
					go game.Start(wg, update, chTG)
					wg.Wait()
				default:
				}
			}
		}
	}
}
