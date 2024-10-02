package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type Message struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

func sendMessage(botToken, receiver, message string) (string, error) {
	api := entrypoint + "sendMessage"

	msgRune := []rune(message)
	var sendingMessage string
	if len(msgRune) > 2048 {
		sendingMessage = string(msgRune[:2048])
	} else {
		sendingMessage = message
	}

	msg := Message{
		ChatID:    receiver,
		Text:      sendingMessage,
		ParseMode: "HTML",
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", api, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	//if resp.StatusCode != http.StatusOK {
	//	return "", errors.New("failed to send message: " + strconv.Itoa(resp.StatusCode))
	//}

	result := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		return "", err
	}

	if result["ok"] != true {
		return "", errors.New("failed to send message: " + result["description"].(string))
	}

	messageId := result["result"].(map[string]interface{})["message_id"].(float64)
	return strconv.FormatFloat(messageId, 'b', 10, 64), nil
}
