package commands

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"os"
	"plugin"
	"strings"
	"unicode"

	"github.com/iCurlmyster/dumb_bot/bot"
	"github.com/iCurlmyster/dumb_bot/bot/parser"
	"github.com/iCurlmyster/dumb_bot/text"
)

// PluginHandler defines what a plugin should accept from the main program
type PluginHandler interface {
	Execute(msg parser.MsgHandler, bot bot.TwitchHandler) error
}

// Plugins holds all loaded plugin objects
type Plugins struct {
	plugins map[string]PluginHandler
}

// LoadPlugins loads plugins from a json file
func LoadPlugins(fileName string) (*Plugins, error) {
	data := map[string]string{}
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}
	plObj := &Plugins{
		plugins: make(map[string]PluginHandler),
	}
	errB := bytes.Buffer{}
	for k, v := range data {
		plug, err := plugin.Open(v)
		if err != nil {
			errB.WriteString(err.Error() + "\n")
			continue
		}
		symPl, err := plug.Lookup("Plugin")
		if err != nil {
			errB.WriteString(err.Error() + "\n")
			continue
		}
		if pl, ok := symPl.(PluginHandler); ok {
			plObj.plugins[k] = pl
		}
	}
	var lastErr error
	if errB.Len() > 0 {
		lastErr = errors.New(string(errB.Bytes()))
	}
	return plObj, lastErr
}

// Execute tries to execute commands in a message
func (p *Plugins) Execute(msg *parser.Msg, bot bot.TwitchHandler) {
	log.Printf("- %s: %s\n", text.Brown(msg.Username()), text.Blue(msg.Message()))
	cmdMsg := strings.TrimSpace(msg.Message())
	if !strings.HasPrefix(cmdMsg, "!") {
		return
	}
	cmdB := bytes.Buffer{}
	for i := 1; i < len(cmdMsg); i++ {
		if unicode.IsSpace(rune(cmdMsg[i])) {
			break
		}
		cmdB.WriteByte(cmdMsg[i])
	}
	cmdMsg = string(cmdB.Bytes())
	if v, ok := p.plugins[cmdMsg]; ok {
		go func() {
			if err := v.Execute(msg.TrimMessagePrefix("!"+cmdMsg), bot); err != nil {
				log.Println(err)
			}
		}()
	}
}
