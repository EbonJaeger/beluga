package beluga

import (
	log2 "log"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"plugin"
	"syscall"

	"github.com/DataDrake/waterlog"
	"github.com/DataDrake/waterlog/format"
	"github.com/DataDrake/waterlog/level"
	"github.com/bwmarrin/discordgo"
)

// Config is our global Config object
var Config RootConfig

// Session is our Discord session
var Session *discordgo.Session

// Log is our WaterLog instance
var Log *waterlog.WaterLog

// Manager is our plugin manager for third-party plugins
var Manager *PluginManager

// ConfigDir is the path to all Beluga-related configs
var ConfigDir string

// Blacklist is the list of users who are blacklisted from using commands
var Blacklist UserBlacklist

// NewBeluga creates a new Beluga instance, and connects to Discord
func NewBeluga(cliFlags Flags) {
	// Initialize logging
	Log = waterlog.New(os.Stdout, "", log2.Ltime)
	Log.SetLevel(level.Info)
	Log.SetFormat(format.Min)

	// Get the path to use for config files
	ConfigDir = cliFlags.ConfigDir
	// Check if we have a path
	if ConfigDir == "" {
		createDefaultConfigDir()
	} else {
		// Make sure the directory exists
		if _, err := os.Stat(ConfigDir); err != nil {
			// Abort on error
			Log.Fatalf("Error while looking for config directory: %s\n", err.Error())
		}
	}

	// Load our config and blacklist
	var readErr error
	if Config, readErr = LoadConfig(); readErr != nil {
		Log.Fatalf("Error while loading config: %s\n", readErr.Error())
	}
	if Blacklist, readErr = LoadBlacklist(); readErr != nil {
		Log.Fatalf("Error while loading or creating user blacklist: %s\n", readErr.Error())
	}

	// Load plugins
	Manager = &PluginManager{
		Plugins: make(map[string]plugin.Symbol),
	}
	if err := Manager.LoadPlugins(); err != nil {
		Log.Fatalf("Error while loading plugins: %s\n", err.Error())
	}

	// Create our Discord client
	Log.Infoln("Creating Discord session")
	s, err := discordgo.New("Bot " + Config.Token)
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

func createDefaultConfigDir() {
	// Get the current user
	var (
		currentUser *user.User
		getUserErr  error
	)
	currentUser, getUserErr = user.Current()
	if getUserErr != nil {
		Log.Fatalf("Unable to get the current user: %s\n", getUserErr.Error())
	}
	// Get the curent user's config directory
	ConfigDir = filepath.Join(currentUser.HomeDir, ".config", "beluga")
	// Check if the config directory exists
	if _, err := os.Stat(ConfigDir); err != nil {
		if os.IsNotExist(err) {
			// Attempt to create the config directory
			if err := os.Mkdir(ConfigDir, 0755); err != nil {
				Log.Fatalf("Unable to create config directory: %s\n", err.Error())
			}
		} else {
			Log.Fatalf("Error while looking for config directory: %s\n", err.Error())
		}
	}
}
