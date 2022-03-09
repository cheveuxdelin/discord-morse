package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/cheveuxdelin/caesar-encoder/morse"
)

var Token string

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// encode
	if strings.HasPrefix(m.Content, "$encode") {

		encodedMessage, error := morse.Encode(strings.TrimPrefix(m.Content, "$encode"))

		if error == nil {
			s.ChannelMessageSend(m.ChannelID, encodedMessage)
		}
	}

	// decode
	if strings.HasPrefix(m.Content, "$decode") {

		encodedMessage, error := morse.Decode(strings.TrimPrefix(m.Content, "$decode"))

		if error == nil {
			s.ChannelMessageSend(m.ChannelID, encodedMessage)
		}
	}

}

func main() {

	// Create a new Discord session using the provided bot token.
	discord, err := discordgo.New("Bot " + Token)

	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	discord.AddHandler(messageCreate)

	discord.Identify.Intents = discordgo.IntentsGuildMessages

	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
}
