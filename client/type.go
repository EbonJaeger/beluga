package client

import (
	log "github.com/DataDrake/waterlog"
	"github.com/EbonJaeger/beluga/handler"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
)

// Client is our Discord handler
type Client struct {
	discord *discordgo.Session
	running bool
	token   string
}

// NewClient will create a new client
func NewClient() *Client {
	return &Client{}
}

func (c *Client) killHandler() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		log.Infoln("Beluga is shutting down")
		c.Close()
		os.Exit(1)
	}()
}

// Close will shut down Beluga and clean everything up
func (c *Client) Close() {
	if !c.running {
		return
	}

	// Shut down services
	c.discord.Close()

	// Mark as stopped
	c.running = false
}

// Start will attempt to start our Discord client
func (c *Client) Start() error {
	// Check of we're already running
	if c.running {
		log.Errorln("Beluga is already running!")
		return nil
	}

	c.running = true
	c.killHandler()

	// Create our Discord client
	d, err := discordgo.New("Bot " + c.token)
	if err != nil {
		return err
	}
	c.discord = d

	// Connect our handlers
	c.discord.AddHandler(handler.OnReady)
	c.discord.AddHandler(handler.OnGuildCreate)
	c.discord.AddHandler(handler.OnMessageCreate)

	// Open Discord websocket
	return c.discord.Open()
}
