package utils

import (
	"math"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func CreateCustomKeyboardByThreeButtons(keys []string) *tg.ReplyKeyboardMarkup {
	if len(keys) == 0 {
		return nil
	}
	retKeyBoard := tg.ReplyKeyboardMarkup{
		Keyboard: make([][]tg.KeyboardButton, 0, int(math.Ceil(float64(len(keys))/3))),
	}
	for i := 0; ; {
		row := make([]tg.KeyboardButton, 0, 3)
		for j := 0; j < 3; j++ {
			if i == len(keys)-1 {
				row = append(row, tg.KeyboardButton{
					Text: keys[i],
				})
				retKeyBoard.Keyboard = append(retKeyBoard.Keyboard, row)
				return &retKeyBoard
			}
			row = append(row, tg.KeyboardButton{
				Text: keys[i],
			})
			i++
		}
		retKeyBoard.Keyboard = append(retKeyBoard.Keyboard, row)
	}
}
