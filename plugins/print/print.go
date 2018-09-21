package main

import (
	"github.com/iCurlmyster/dumb_bot/bot"
	"github.com/iCurlmyster/dumb_bot/bot/parser"
)

type printType struct{}

func (p printType) Execute(msg parser.MsgHandler, bot bot.TwitchHandler) error {
	bot.WriteMessage("your message: " + msg.Message())
	return nil
}

// Plugin the main object for this plugin
var Plugin printType
