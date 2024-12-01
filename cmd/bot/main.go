package main

import (
	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/indexcoder/oauth2-authorization/pkg/pocket"
	"github.com/indexcoder/oauth2-authorization/pkg/repository"
	"github.com/indexcoder/oauth2-authorization/pkg/repository/boltdb"
	"github.com/indexcoder/oauth2-authorization/pkg/server"
	"github.com/indexcoder/oauth2-authorization/pkg/telegram"
	"log"
)

var userStates = make(map[int64]string)

func main() {
	// Создаем бота
	bot, err := tgbotapi.NewBotAPI("7543227307:AAGYpAkSZofJDLv5SIKIm2nETeLO0cxwTzw")
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = true

	pocketClient, err := pocket.NewClient("111953-4f0a374f95750a63cab5a5e")
	if err != nil {
		log.Fatal(err)
	}

	db, err := initDb()
	if err != nil {
		log.Fatal(err)
	}

	tokenRepository := boltdb.NewTokenRepository(db)

	telegramBot := telegram.NewBot(bot, pocketClient, tokenRepository, "http://localhost/")

	authorizationServer := server.NewAuthorizationServer(pocketClient, tokenRepository, "https://t.me/LfScoutBot")

	go func() {
		if err := telegramBot.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := authorizationServer.Start(); err != nil {
		log.Fatal(err)
	}

}

func initDb() (*bolt.DB, error) {
	db, err := bolt.Open("bot.db", 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.AccessToken))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(repository.RequestTokens))
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return db, nil
}
