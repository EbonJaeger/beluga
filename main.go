package beluga

import (
	log2 "log"
	"os"
	"os/signal"
	"plugin"
	"syscall"

	"github.com/DataDrake/waterlog"
	"github.com/DataDrake/waterlog/format"
	"github.com/DataDrake/waterlog/level"
	"github.com/EbonJaeger/beluga/config"
	"github.com/bwmarrin/discordgo"
)

// Session is our Discord session
var Session *discordgo.Session

// Log is our WaterLog instance
var Log *waterlog.WaterLog

// PluginManager is our plugin manager for third-party plugins
var PluginManager *BelugaPluginManager

// NewBeluga creates a new Beluga instance, and connects to Discord
func NewBeluga() {
	// Initialize logging
	Log = waterlog.New(os.Stdout, "", log2.Ltime)
	Log.SetLevel(level.Info)
	Log.SetFormat(format.Min)

	// Load plugins
	PluginManager = &BelugaPluginManager{
		Plugins: make(map[string]plugin.Symbol),
	}
	if err := PluginManager.LoadPlugins(); err != nil {
		Log.Fatalf("Error while loading plugins: %s\n", err.Error())
	}

	// Create our Discord client
	Log.Infoln("Creating Discord session")
	s, err := discordgo.New("Bot " + config.Conf.Token)
	if err != nil {
		Log.Fatalf("Unable to initialize discordgo: %s\n", err.Error())
	}
	Session = s

	// Connect our handlers
	Session.AddHandler(OnReady)
	Session.AddHandler(OnGuildCreate)
	Session.AddHandler(OnMessageCreate)

	// Open Discord websocket
	Log.Infoln("Connecting to Discord websocket")
	if err := Session.Open(); err != nil {
		Log.Fatalf("Unable to connect to Discord websocket: %s\n", err.Error())
	}

	// Wait until told to close
	Log.Goodln("Connected to Discord! Press CTRL+C to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	// Get to a new line
	Log.Println("")

	// Close the Discord session on close
	if err = Session.Close(); err != nil {
		Log.Fatalf("Error while closing Discord connection: %s\n", err.Error())
	}
	Log.Goodln("Beluga shut down successfully!")
}
