package parser

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"unicode"
)

// MsgType represents Messages types that can be returned from twitch irc
type MsgType int

const (
	// PrivMsg represents a private message type
	PrivMsg = MsgType(iota)
	// Notice represents a Notice type from server
	Notice
	// Ping represents a ping message type
	Ping
)

var (
	pingB = []byte("PING")
	tmiB  = []byte(":tmi")
)

// Msg represents
type Msg struct {
	username string
	message  string
	msgType  MsgType
}

// MsgHandler defines what a plugin can get from a Msg object
type MsgHandler interface {
	Username() string
	Message() string
}

// TrimMessagePrefix used mainly to shave off command messages when sending message to plugin
func (m *Msg) TrimMessagePrefix(pre string) *Msg {
	return &Msg{
		username: m.username,
		message:  strings.TrimPrefix(strings.TrimSpace(m.message), pre),
		msgType:  m.msgType,
	}
}

// Username returns the Username of the Message
func (m *Msg) Username() string {
	return m.username
}

// Message returns the user's message
func (m *Msg) Message() string {
	return m.message
}

// Type returns the message type
func (m *Msg) Type() MsgType {
	return m.msgType
}

// TwitchMessage parses a message from twitch and returns a Msg object
func TwitchMessage(r io.Reader) *Msg {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		log.Println(err)
		return nil
	}
	front := b[:4]
	if bytes.Compare(pingB, front) == 0 {
		return &Msg{
			msgType: Ping,
		}
	} else if bytes.Compare(tmiB, front) == 0 {
		return handleTmiMsg(b)
	}
	return handleUserMsg(b)
}

func handleTmiMsg(b []byte) *Msg {
	tail := b[15:]
	if unicode.IsDigit(rune(tail[0])) {
		code, err := strconv.Atoi(string(tail[:3]))
		if err != nil {
			log.Println(err)
			return nil
		}
		if code == 421 {
			log.Println("Command was not recognized")
			return nil
		}
	} else {
		msg := &Msg{}
		word, index := getUntilDelim(tail, 0, ' ')
		if string(word) == "NOTICE" {
			un, index := getUntilDelim(tail, index+1, ' ')
			// skip 2 characters to account for " :"
			m, index := getUntilDelim(tail, index+2, '\n')
			msg.username = string(un)
			msg.message = string(m)
			msg.msgType = Notice
			return msg
		}
	}
	return nil
}

func handleUserMsg(b []byte) *Msg {
	msg := &Msg{}
	// start at 1 to skip the ":"
	un, index := getUntilDelim(b, 1, '!')
	index = bytes.IndexByte(b, ' ')
	if index == -1 {
		return nil
	}
	command, index := getUntilDelim(b, index+1, ' ')
	if string(command) == "PRIVMSG" {
		_, index = getUntilDelim(b, index+1, ':')
		if index == -1 {
			return nil
		}
		m, _ := getUntilDelim(b, index+1, '\n')
		msg.username = string(un)
		msg.message = string(m)
		msg.msgType = PrivMsg
		return msg
	}
	return nil
}

func getUntilDelim(b []byte, index int, delim byte) ([]byte, int) {
	bLen := len(b)
	tmp := make([]byte, 0)
	if index >= bLen {
		return tmp, 0
	}
	for i := index; ; i++ {
		if i >= bLen || b[i] == delim {
			return tmp, i
		}
		tmp = append(tmp, b[i])
	}
}
