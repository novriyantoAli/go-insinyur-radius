package telegram

import (
	"context"
	"insinyur-radius/domain"
	"log"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// NewHandler ...
func NewHandler(ucm domain.MessageUsecase) {
	var mtx sync.Mutex

LoopReconnect:
	for true {
		bot, err := tgbotapi.NewBotAPI(viper.GetString("telegram.token"))
		if err != nil {
			logrus.Error(err)
			time.Sleep(5 * time.Second)
			
			continue LoopReconnect
		}

		bot.Debug = true

		log.Printf("Authorized on account %s", bot.Self.UserName)

		for true {
			mtx.Lock()

			recParams := "no"
			spec := domain.Message{
				ID:        nil,
				ChatID:    nil,
				MessageID: nil,
				Received:  &recParams,
				Message:   nil,
				CreatedAt: nil,
			}

			res, err := ucm.Find(context.Background(), spec)
			if err != nil {
				logrus.Error(err)
			}

			if len(res) > 0 {
				for _, message := range res {
					msg := tgbotapi.NewMessage(*message.ChatID, *message.Message)
					msg.ReplyToMessageID = int(*message.MessageID)
					_, err := bot.Send(msg)
					if err != nil {
						logrus.Error(err)
					} else {
						received := "yes"
						updated := domain.Message{}
						updated.ID = message.ID
						updated.ChatID = message.ChatID
						updated.MessageID = message.MessageID
						updated.Received = &received

						err := ucm.Update(context.Background(), updated)
						if err != nil {
							logrus.Error(err)
						}
					}
				}
			}
			mtx.Unlock()
			time.Sleep(5 * time.Second)
		}
	}
}
