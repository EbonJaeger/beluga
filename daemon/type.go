package daemon

import (
	log "github.com/DataDrake/waterlog"
	"github.com/EbonJaeger/beluga/api/v1"
	"github.com/EbonJaeger/beluga/config"
	"github.com/EbonJaeger/beluga/handler"
	"github.com/EbonJaeger/beluga/plugins"
	"github.com/bwmarrin/discordgo"
	"github.com/coreos/go-systemd/daemon"
	"os"
	"os/signal"
	"plugin"
	"syscall"
)

// Server is our Discord handler
type Server struct {
	api           *v1.Listener
	discord       *discordgo.Session
	pluginManager *plugins.PluginManager
	running       bool
}

// NewServer will create a new daemon
func NewServer() *Server {
	return &Server{}
}

func (s *Server) killHandler() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		log.Infoln("Beluga is shutting down")
		s.Close()
		os.Exit(1)
	}()
}

// Close will shut down Beluga and clean everything up
func (s *Server) Close() {
	if !s.running {
		return
	}

	// Shut down services
	s.discord.Close()
	s.api.Close()

	// Mark as stopped
	s.running = false
}

// Bind will set up a listener on our Unix socket
func (s *Server) Bind() error {
	// Set up the API
	api, err := v1.NewListener()
	if err != nil {
		return err
	}
	s.api = api
	return s.api.Bind()
}

// Start will attempt to start our Discord client
func (s *Server) Start() error {
	// Check of we're already running
	if s.running {
		log.Errorln("Beluga is already running!")
		return nil
	}

	s.running = true
	s.killHandler()

	// Load plugins
	pm := plugins.NewManager()
	pm.Plugins = make(map[string]plugin.Symbol)
	if err := pm.LoadPlugins(); err != nil {
		return err
	}
	s.pluginManager = pm

	// Initilize handler funcs struct
	var funcs = &handler.Funcs{
		PluginManager: pm,
	}

	// Create our Discord client
	d, err := discordgo.New("Bot " + config.Conf.Token)
	if err != nil {
		return err
	}

	// Connect our handlers
	d.AddHandler(funcs.OnReady)
	d.AddHandler(funcs.OnGuildCreate)
	d.AddHandler(funcs.OnMessageCreate)

	s.discord = d
	s.pluginManager.Session = d

	// Open Discord websocket
	if err := s.discord.Open(); err != nil {
		return err
	}

	// Notify Systemd
	if s.api.SystemdEnabled {
		ok, err := daemon.SdNotify(false, daemon.SdNotifyReady)
		if err != nil {
			log.Errorf("Failed to notify Systemd: %s\n", err.Error())
			return err
		}
		if !ok {
			log.Warnln("SdNotify failed because of missing environment variable")
		} else {
			log.Goodln("SdNotify successful")
		}
	}

	// Start the API server
	return s.api.Start()
}
