package telegram

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
)

//pics for game avatar
const (
	GameStartPic = "internal/client/telegram/base_images/game_begin.jpeg"
	GameEndPic   = "internal/client/telegram/base_images/game_end.jpeg"
)

type Client interface {
	SendMessage(int64, string) error
	SendMessageWithKeyboard(chatID int64, message string, keyboard tg.ReplyKeyboardMarkup) error
	SetAvatar(chatID int64, avatarPath string) error
	SendMessageAndReturnMSG(chatID int64, message string) (tg.Message, error)
}

type client struct {
	bot        *tg.BotAPI
	updateChan tg.UpdatesChannel
}

func NewTelegramClient(token string) (Client, tg.UpdatesChannel, error) {
	bot, err := tg.NewBotAPI(token)
	if err != nil {
		return nil, nil, err
	}

	bot.Debug = true

	cfg := tg.NewUpdate(0)
	cfg.Timeout = 60

	chanTG, err := bot.GetUpdatesChan(cfg)
	if err != nil {
		return nil, nil, err
	}
	return &client{bot, chanTG}, chanTG, nil
}

func (c *client) SendMessage(chatID int64, message string) error {
	msg := tg.NewMessage(chatID, message)

	_, err := c.bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func (c *client) SendMessageAndReturnMSG(chatID int64, message string) (tg.Message, error) {
	msg := tg.NewMessage(chatID, message)

	mSend, err := c.bot.Send(msg)
	if err != nil {
		return tg.Message{}, err
	}
	return mSend, nil
}

func (c *client) SendMessageWithKeyboard(chatID int64, message string, keyboard tg.ReplyKeyboardMarkup) error {
	msg := tg.NewMessage(chatID, message)

	keyboard.OneTimeKeyboard = true

	msg.ReplyMarkup = keyboard

	_, err := c.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) SetAvatar(chatID int64, avatarPath string) error {
	f, err := ioutil.ReadFile(avatarPath)
	if err != nil {
		return err
	}

	cfg := tg.SetChatPhotoConfig{
		BaseFile: tg.BaseFile{
			BaseChat: tg.BaseChat{ChatID: chatID},
			File: tg.FileBytes{
				Name:  "gameStartPic",
				Bytes: f,
			},
		},
	}
	resp, err := c.bot.SetChatPhoto(cfg)
	if err != nil || resp.Ok != true {
		fmt.Println(string(resp.Result))
		return err
	}
	return nil
}
