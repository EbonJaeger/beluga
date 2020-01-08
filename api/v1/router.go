package v1

import (
	"errors"
	log "github.com/DataDrake/waterlog"
	"github.com/EbonJaeger/beluga/config"
	"github.com/coreos/go-systemd/activation"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"net"
	"net/http"
	"os"
	"time"
)

// Listener listens on a Unix socket for connections from
// an authenticated user
type Listener struct {
	srv            *fasthttp.Server
	router         *router.Router
	socket         net.Listener
	SystemdEnabled bool
	timeStarted    time.Time
}

// NewListener creates a new unbound Server, and initializes
// it
func NewListener() (api *Listener, err error) {
	r := router.New()
	api = &Listener{
		srv: &fasthttp.Server{
			Handler: r.Handler,
		},
		router:      r,
		timeStarted: time.Now().UTC(),
	}

	// Define the API endpoints
	r.DELETE("/api/v1/daemon", api.StopDaemon)

	return api, nil
}

// Bind will set up the listner on the Unix socket
func (api *Listener) Bind() error {
	var listener net.Listener
	log.Infoln("Here???")

	// Check for Systemd
	if val, present := os.LookupEnv("LISTEN_FDS"); present {
		log.Debugf("Env 'LISTEN_FDS' is set: %v\n", val)
		listeners, err := activation.Listeners()
		if err != nil {
			return err
		}
		if len(listeners) != 1 {
			return errors.New("Expected a single Unix socket")
		}

		listener = listeners[0]
		if unix, ok := listener.(*net.UnixListener); ok {
			unix.SetUnlinkOnClose(false)
		} else {
			return errors.New("Not a Unix socket")
		}

		api.SystemdEnabled = true
	} else {
		l, err := net.Listen("unix", config.Conf.Socket)
		if err != nil {
			return err
		}
		listener = l
	}

	uid := os.Getuid()
	gid := os.Getgid()
	if !api.SystemdEnabled {
		// Make sure we own the socket
		if err := os.Chown(config.Conf.Socket, uid, gid); err != nil {
			return err
		}
		if err := os.Chmod(config.Conf.Socket, 0660); err != nil {
			return err
		}
	}

	api.socket = listener
	return nil
}

// Start will serve on the Unix socket until it is killed
func (api *Listener) Start() error {
	if api.socket == nil {
		return errors.New("Cannot serve without a bound socket")
	}

	if err := api.srv.Serve(api.socket); err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Close will shut down the server and clean up the socket
func (api *Listener) Close() {
	api.srv.Shutdown()
	if !api.SystemdEnabled {
		os.Remove(config.Conf.Socket)
	}
}
