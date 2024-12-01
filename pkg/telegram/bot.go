package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/indexcoder/oauth2-authorization/pkg/pocket"
	"github.com/indexcoder/oauth2-authorization/pkg/repository"
	"log"
)

type Bot struct {
	bot             *tgbotapi.BotAPI
	pocketClient    *pocket.Client
	tokenRepository repository.TokenRepository
	redirectURL     string
}

func NewBot(bot *tgbotapi.BotAPI, pocketClient *pocket.Client, tr repository.TokenRepository, redirectUrl string) *Bot {
	return &Bot{bot: bot, pocketClient: pocketClient, tokenRepository: tr, redirectURL: redirectUrl}
}

func (b *Bot) Start() error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)
	updates := b.initUpdatesChannel()
	b.handleUpdates(updates)
	return nil
}

func (b *Bot) initUpdatesChannel() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return b.bot.GetUpdatesChan(u)
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message != nil {

			if update.Message.IsCommand() {
				err := b.handleCommand(update.Message)
				if err != nil {
					log.Println(err)
				}
				continue
			}

			b.handleMessage(update.Message)
		}
	}
}
