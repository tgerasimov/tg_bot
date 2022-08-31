package hook

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type hookFunc func(update tgbotapi.Update)

//Hook hook for checking is user have permission to access command
func Hook(ctx context.Context, f hookFunc) {

}
