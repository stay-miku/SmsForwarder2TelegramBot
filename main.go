package main

import (
	"log"
	"net/http"
)

func main() {
	log.Printf("Config loaded. Bot Token: %s, Receiver: %s, Port: %s\n", config.BotToken, config.Receiver, config.Port)

	http.HandleFunc("/post", process)

	if config.Https.Enable {
		log.Printf("Enable HTTPS")
		err := http.ListenAndServeTLS(":"+config.Port, config.Https.Cert, config.Https.Key, nil)
		if err != nil {
			panic(err)
		}
		return
	} else {
		err := http.ListenAndServe(":"+config.Port, nil)
		if err != nil {
			panic(err)
		}
	}
}
