package daemon

import (
	log "github.com/DataDrake/waterlog"
	"github.com/EbonJaeger/beluga/api/v1"
	"github.com/EbonJaeger/beluga/config"
	"github.com/EbonJaeger/beluga/handler"
	"github.com/bwmarrin/discordgo"
	"github.com/coreos/go-systemd/daemon"
	"os"
	"os/signal"
	"syscall"
)

// Server is our Discord handler
type Server struct {
	api     *v1.Listener
	discord *discordgo.Session
	running bool
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

	// Create our Discord client
	d, err := discordgo.New("Bot " + config.Conf.Token)
	if err != nil {
		return err
	}

	// Connect our handlers
	d.AddHandler(handler.OnReady)
	d.AddHandler(handler.OnGuildCreate)
	d.AddHandler(handler.OnMessageCreate)

	s.discord = d

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
