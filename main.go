package main

import (
	"log"
	"net/http"
)

func main() {
	log.Printf("Config loaded. Bot Token: %s, Receiver: %s, Port: %s\n", config.BotToken, config.Receiver, config.Port)

	http.HandleFunc("/post", process)

	err := http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		panic(err)
	}
}
