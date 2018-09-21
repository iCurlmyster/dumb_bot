package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/iCurlmyster/dumb_bot/bot"
	"github.com/iCurlmyster/dumb_bot/bot/commands"
	"github.com/iCurlmyster/dumb_bot/bot/parser"
	"github.com/iCurlmyster/dumb_bot/config"
)

func main() {
	c, err := config.Configuration("./resources/config.json")
	if err != nil {
		panic(err)
	}
	tw := bot.NewTwitchClient(c)
	fmt.Println("Connected!")
	defer tw.Close()

	cl := make(chan bool)
	msgAlert := make(chan *parser.Msg)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	pl, err := commands.LoadPlugins("./resources/plugins.json")
	if err != nil {
		log.Println(err)
	}

	go handleListenForMessages(msgAlert, tw, cl)()

	for {
		select {
		case msg := <-msgAlert:
			if msg.Type() == parser.PrivMsg {
				pl.Execute(msg, tw)
			} else if msg.Type() == parser.Ping {
				if err := tw.Pong(); err != nil {
					log.Println(err)
				}
			}
		case <-quit:
			log.Println("calling close...")
			tw.Close()
			cl <- true
			return
		}
	}
}

func handleListenForMessages(msgAlert chan *parser.Msg, tw *bot.Twitch, cl chan bool) func() {
	return func() {
		for {
			select {
			case <-cl:
				return
			default:
				b, err := tw.ListenForMessage()
				if err != nil {
					log.Fatalf("listen error: %v", err)
				}
				msg := parser.TwitchMessage(b)
				if msg != nil {
					msgAlert <- msg
				}
			}
		}
	}
}