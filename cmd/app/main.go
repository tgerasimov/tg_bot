package main

import (
	"bot_tg/internal/config"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"bot_tg/internal/client/telegram"
	user_list "bot_tg/internal/client/user-list"
	"bot_tg/internal/commands"
)

func main() {

	ctx := context.Background()

	cfg, ok := config.GetConfigFromCtx(config.ParseConfig(ctx))
	if !ok {
		log.Fatal("can't parse config")
	}

	token := cfg.GetValue(config.Telegram_token)
	if token == "" {
		log.Fatal("can't parse token")
	}
	tgClient, chanTG, err := telegram.NewTelegramClient(token)
	if err != nil {
		log.Fatal(err)
	}

	userList := user_list.NewUserList()

	commands.RegisterAllCommands(ctx, chanTG, tgClient, userList)

	graceful() //TODO() connect everything through the graceful
}

func graceful() {
	exit := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	for {
		select {
		case <-exit:
			log.Println("Shutting down")
			return
		default:
		}
	}
}
