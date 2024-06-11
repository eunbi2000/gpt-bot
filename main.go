package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/jmorganca/ollama/api"
)

// Variables used for command line parameters
const Token = ""

var client *discordgo.Session
var apiClient *api.Client

func main() {
	session, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Error creating Discord session:", err)
		return
	}
	client = session
	client.AddHandler(message)
	apiClient = api.NewClient()
	fmt.Println("Bot is online")
	defer client.Close()
	if err = client.Open(); err != nil {
		fmt.Println("Error opening connection:", err)
		return
	}

	scall := make(chan os.Signal, 1)
	signal.Notify(scall, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSEGV, syscall.SIGHUP)
	<-scall
	fmt.Println("\nBot shutting down.")
}

func message(bot *discordgo.Session, message *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if message.Author.ID == bot.State.User.ID {
		return
	}

	var latest api.GenerateResponse
	var output string
	generateContext := []int{} // TODO: Get conversation context

	request := api.GenerateRequest{Model: "llama3", Prompt: message.Content, Context: generateContext}
	fn := func(response api.GenerateResponse) error {
		latest = response
		output += response.Response
		fmt.Print(response.Response)
		return nil
	}

	// Log request details
	fmt.Printf("Sending request to LLaMA: %+v\n", request)

	err := apiClient.Generate(context.Background(), &request, fn)
	if err != nil {
		fmt.Println("Error generating response from LLaMA:", err)
		bot.ChannelMessageSend(message.ChannelID, "There seems to be a problem, please cry to the developer")
		return
	}
	fmt.Println(" \n ")
	bot.ChannelMessageSend(message.ChannelID, output)
	latest.Summary()
}
