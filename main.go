package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"

	"github.com/JoshStrobl/trunk"
)

var token string

func init() {
	flag.StringVar(&token, "t", "", "Discord bot token")
	flag.Parse()
}

func main() {
	if token == "" { // No bot token provided
		trunk.LogFatal("No bot token was provided!")
	}

	discord, err := discordgo.New("Bot " + token)
	if err != nil { // Error occurred while creating Discord session
		trunk.LogFatal("Unable to create Discord session: " + err.Error())
	}

	/* Add our handlers */
	discord.AddHandler(onReady)
	discord.AddHandler(onMessageCreate)

	/* Open websocket and start listening */
	err = discord.Open()
	if err != nil { // Error while opening Discord session
		trunk.LogFatal("Unable to start Discord session: " + err.Error())
	}

	// Wait until told to end
	trunk.LogSuccess("Beluga is now running! Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close the Discord session
	discord.Close()
}

// onReady is called when the bot receives the "ready" event from Discord.
func onReady(session *discordgo.Session, event *discordgo.Ready) {
	session.UpdateStatus(0, "Moderation")
}

func onGuildCreate(session *discordgo.Session, event *discordgo.GuildCreate) {
	if event.Guild.Unavailable { // Joined Guild is unavailable
		trunk.LogWarn("Attempted to join a guild but it was unavailable")
		return
	}

	for _, channel := range event.Guild.Channels { // Iterate through the Guild's channels
		if channel.ID == event.Guild.ID { // I have no idea when this is true
			_, _ = session.ChannelMessageSend(channel.ID, "Beluga is ready!")
			return
		}
	}
}

// onMessageCreate is called whenever a message is created in a channel that
// the bot has access to.
func onMessageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID { // If the message was sent by this bot
		return
	}

	trunk.LogInfo(fmt.Sprintf("Message received in channel '%s': %s", message.ChannelID, message.Content))
}
