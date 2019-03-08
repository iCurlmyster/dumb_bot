package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/iCurlmyster/dumb_bot/bot"
	"github.com/iCurlmyster/dumb_bot/bot/commands"
	"github.com/iCurlmyster/dumb_bot/bot/parser"
	"github.com/iCurlmyster/dumb_bot/config"
	"github.com/iCurlmyster/dumb_bot/text"
)

func main() {
	c, err := config.Configuration("./resources/config.json")
	if err != nil {
		panic(err)
	}
	tw := bot.NewTwitchClient(c)
	log.Println(text.Cyan("Connected!"))
	defer tw.Close()

	msgAlert := make(chan *parser.Msg)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	defer close(msgAlert)

	pl, err := commands.LoadPlugins("./resources/plugins.json")
	if err != nil {
		log.Println(err)
	}

	go handleListenForMessages(msgAlert, tw)()

	for {
		select {
		case msg, ok := <-msgAlert:
			if !ok {
				break
			}
			if msg.Type() == parser.PrivMsg {
				pl.Execute(msg, tw)
			} else if msg.Type() == parser.Ping {
				if err := tw.Pong(); err != nil {
					log.Println(err)
				}
			}
		case <-quit:
			log.Println(text.Cyan("Closing connection"))
			tw.Close()
			break
		}
	}
}

func handleListenForMessages(msgAlert chan *parser.Msg, tw *bot.Twitch) func() {
	return func() {
		for {
			b, err := tw.ListenForMessage()
			if err != nil {
				log.Fatalf(text.Red("read error: %v\n"), err)
				close(msgAlert)
				break
			}
			msg := parser.TwitchMessage(b)
			if msg != nil {
				msgAlert <- msg
			}
		}
	}
}
