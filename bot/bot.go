package bot

import (
	"bytes"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/iCurlmyster/dumb_bot/config"
)

// TwitchHandler defines what a plugin can expect to be able to interact with
type TwitchHandler interface {
	WriteMessage(msg string) error
}

// Twitch holds fields required to interact with twitch api
type Twitch struct {
	client        *websocket.Conn
	configuration *config.Config
}

// NewTwitchClient generates a new Twitch object with the passed in configurations
func NewTwitchClient(c *config.Config) *Twitch {
	tw := &Twitch{
		configuration: c,
	}
	conn, resp, err := websocket.DefaultDialer.Dial("wss://irc-ws.chat.twitch.tv:443", nil)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	printReader(resp.Body)
	tw.client = conn
	passStr := fmt.Sprintf("PASS %s", tw.configuration.OAuth)
	tw.client.WriteMessage(websocket.TextMessage, []byte(passStr))
	nickStr := fmt.Sprintf("NICK %s", tw.configuration.UserName)
	tw.client.WriteMessage(websocket.TextMessage, []byte(nickStr))
	chanStr := fmt.Sprintf("JOIN #%s", tw.configuration.Channel)
	tw.client.WriteMessage(websocket.TextMessage, []byte(chanStr))
	return tw
}

// WriteMessage writes the given string as a PRIVMSG
func (tw *Twitch) WriteMessage(msg string) error {
	msgStr := fmt.Sprintf("PRIVMSG #%s :%s", tw.configuration.Channel, msg)
	return tw.client.WriteMessage(websocket.TextMessage, []byte(msgStr))
}

// ListenForMessage grabs the latest message from the websocket as a Buffer object
func (tw *Twitch) ListenForMessage() (*bytes.Buffer, error) {
	_, r, err := tw.client.NextReader()
	if err != nil {
		return nil, err
	}
	return printReader(r)
}

// Pong sends a pong message to the server
func (tw *Twitch) Pong() error {
	return tw.client.WriteMessage(websocket.TextMessage, []byte("PONG :tmi.twitch.tv"))
}

// Close handles closing internal client objects
func (tw *Twitch) Close() {
	tw.client.Close()
}
