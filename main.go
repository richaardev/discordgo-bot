package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main()  {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error while loading .env file!")
	}

	token := os.Getenv("TOKEN")

	discord, derr := dgo.New("Bot " + token)
	if derr != nil {
		panic(derr)
	}
	
	discord.AddHandlerOnce(func (s *dgo.Session, m *dgo.Ready) {
		fmt.Println("The super", s.State.User.Username, "is now online!")
	})
	discord.AddHandler(onMessage)

	
	if err = discord.Open(); err != nil {
		fmt.Printf("Opening the session failed: \"%s\".\n", err.Error())
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")


	WaitForInterrupt()
	fmt.Println("Shutting down.")
	if err := discord.Close(); err != nil {
		fmt.Printf("Closing the session failed. \"%s\". Ignoring.\n", err.Error())
	}
}
func onMessage(s *dgo.Session, m *dgo.MessageCreate) {

	if m.Author.Bot { return }
	
	prefix := "!"
	
	args := strings.Split(strings.TrimPrefix(m.Content, prefix), " ")
	command := args[0]
	channel, cerr := s.State.GuildChannel(m.GuildID, m.ChannelID)
	if cerr != nil  {
		log.Fatalln("Cannot find the channel!")
	}
	if command == "ping" {
		s.ChannelMessageSend(channel.ID, "Pong!")
	}
}


func WaitForInterrupt() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}