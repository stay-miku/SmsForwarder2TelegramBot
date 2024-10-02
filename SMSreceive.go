package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"time"
)

var logger = log.New(os.Stdout, "SMS:", log.LstdFlags)

func messageGenerate(sender, content string) string {
	re := regexp.MustCompile(`\d{4,}`)
	content = re.ReplaceAllString(content, "<code>$0</code>")

	return "SMS from: <code>" + sender + "</code>\n----------------------------------------\n" + content
}

func validateSign(sign, timestamp string) bool {
	currenTime := time.Now().UnixMilli()
	intTimestamp, err := strconv.ParseInt(timestamp, 10, 64)

	if err != nil {
		return false
	}

	// check time
	if currenTime-intTimestamp > 300000 || currenTime-intTimestamp < -300000 {
		return false
	}

	secretWithTimestamp := []byte(timestamp + "\n" + config.Secret)
	secret := []byte(config.Secret)

	h := hmac.New(sha256.New, secret)
	h.Write(secretWithTimestamp)
	hmacSign := h.Sum(nil)

	b64Sign := base64.StdEncoding.EncodeToString(hmacSign)
	parseSign := url.QueryEscape(b64Sign)

	if parseSign != sign {
		return false
	}
	return true
}

func process(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// get form and client addr
	//client := r.RemoteAddr
	client := r.Header.Get("CF-Connecting-IP")
	if client == "" {
		client = r.RemoteAddr
	}
	sender := r.Form.Get("from")
	content := r.Form.Get("content")
	timestamp := r.Form.Get("timestamp")
	sign := r.Form.Get("sign")

	// check form
	if sender == "" || content == "" || timestamp == "" || sign == "" {
		http.Error(w, "Missing parameters", http.StatusBadRequest)
		logger.Printf("Missing parameters. Client: %s\n", client)
		return
	}

	// check sign
	if !validateSign(sign, timestamp) {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		logger.Printf("Authentication failed. Client: %s\n", client)
		return
	}

	logger.Printf("Request Received. Client: %s, Sender: %s, Content: %s, Timestamp: %s, Sign: %s\n", client, sender, content, timestamp, sign)

	// send message via tg bot
	result, err := sendMessage(config.BotToken, config.Receiver, messageGenerate(sender, content))
	if err != nil {
		http.Error(w, "Error sending message", http.StatusInternalServerError)
		logger.Printf("Error sending message. Client: %s, Error: %s\n", client, err)
		return
	} else {
		logger.Printf("Message sent successfully. Client: %s, Message ID: %s\n", client, result)
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("send successfully, message id: " + result))
	if err != nil {
		logger.Printf("Error writing response. Client: %s\n", client)
		return
	}
}
