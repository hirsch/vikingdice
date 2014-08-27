package main

import (
	"github.com/hirsch/ircudf"
	"github.com/hirsch/conf"
	"log"
	"math/rand"
	"strings"
)

var (
	chann = ""
	admin = ""
	
	add = make(map[string]string)
)

func main() {
	cfg, err := conf.Open("server.cfg"); check(err)
	addr, err := cfg.Read("server", "address"); check(err)
	pass, err := cfg.Read("server", "password"); check(err)
	nick, err := cfg.Read("server", "nickname"); check(err)
	chann, err = cfg.Read("server", "channel"); check(err)
	admin, err = cfg.Read("server", "admin"); check(err)
	
	admin = strings.ToLower(admin)
	
	ircudf.Debug = true
	ircudf.HandleReply(reply)
	ircudf.HandlePrivmsg(privmsg)
	
	server := ircudf.Create(addr, nick, nick, nick, pass)
	err = server.Connect(); check(err)
}

func reply(server *ircudf.Server, number, name, reply string) {
	if number == "376" {
		server.Join(chann)
	}
}

func privmsg(server *ircudf.Server, channel, user, message string) {
	
	if message == "!get" && strings.ToLower(user) == admin {
		if len(add) < 1 {
			server.Privmsg(chann, "no entries yet")
			return
		}
		i := 0
		rand := rand.Intn(len(add))
		
		for name, entry := range add {
			if i == rand {
				server.Privmsg(chann, "Selected " + name + ": " + entry)
				break
			}
			i++
		}
		add = make(map[string]string)
		return
	}
	
	if len(message) > 5 {
		if message[:5] == "!add " {
			add[user] = message[5:]
		}
	}
}


func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
