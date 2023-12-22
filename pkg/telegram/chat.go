package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

type telegramMessage struct {
	ChatId string `json:"chat_id"`
	Text   string `json:"text"`
}

type Telegram struct {
	Token string
}

func (telegram Telegram) getUrl() string {
	return fmt.Sprintf("https://api.telegram.org/bot%s", telegram.Token)
}

func (telegram Telegram) SendMessage(chatId string, text string) error {
	url := fmt.Sprintf("%s/sendMessage", telegram.getUrl())
	message := telegramMessage{
		ChatId: chatId,
		Text:   text,
	}
	body, _ := json.Marshal(message)
	response, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Err(err).Msg("Error posting message to telegram")
		return err
	}
	defer response.Body.Close()
	body, err = io.ReadAll(response.Body)
	if err != nil {
		log.Err(err).Msg("Error decoding response body from sending message to telegram")
		return err
	}
	log.Info().Msgf("Message '%s' was sent, with response '%s'", text, body)
	return nil
}
