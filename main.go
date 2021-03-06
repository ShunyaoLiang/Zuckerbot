package main

import (
	"math/rand"
	"time"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters

var userBlacklist = []string {
	""
}

var channelBlacklist = []string {
        ""
}

func init() {
	// Add nodes representing the beginning and ending of sentences
	markov = append(markov, node{"FRONT", make([]*node, 0)})
	markov = append(markov, node{"BACK", make([]*node, 0)})

	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot Token")
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(onMessage)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// See if sender is blacklisted
	for _, v := range userBlacklist {
		if v == m.Author.ID {
			return
		}
	}
	// See if sender is bot
	if m.Author.Bot {
		return
	}
	// See if message is from whitelisted channel
	allowed := false
	for _, v := range channelWhitelist {
		if m.ChannelID == v {
			allowed = true
		}
	}
	if !allowed {
		return
	}

	// Check for #ask-zucc event
	if m.ChannelID == "468402235412578314" && (strings.Contains(m.Message.Content, "Zuckerbot") || strings.Contains(m.Message.Content, "zuckerbot")) {
		fmt.Println("QUESTION RECEIVED MORTAL")
		s.ChannelMessageSend(m.ChannelID, generate())
		return
	}

	// Clean and prepare message for interpretation
	message := strings.ToLower(m.ContentWithMentionsReplaced())
	line := strings.Split(message, "\n")
	data := make([]string, 0)
	for _, v := range line {
		data = append(data, strings.Split(v, " ")...)
	}
	fmt.Println(data)

	getWord := func(str string) *node {
		for i := range markov {
			if markov[i].word == str {
				return &markov[i]
			}
		}

		markov = append(markov, node{str, make([]*node, 0)})
		for i := range markov {
			if markov[i].word == str {
				return &markov[i]
			}
		}

		panic("If you are reading this message, the world is ending")
	}

	// Interpret the data
	getWord("FRONT").addLink(getWord(data[0]))
	for i := range data[:len(data)-1] {
		getWord(data[i]).addLink(getWord(data[i+1]))
	}
	getWord(data[len(data)-1]).addLink(getWord("BACK"))

	// Randomly decide if Zuckerbot should speak
	if rand.Intn(4) == 4 {
		s.ChannelMessageSend(m.ChannelID, generate())
	}
}
