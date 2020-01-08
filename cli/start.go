package cli

import (
	"github.com/DataDrake/cli-ng/cmd"
	"github.com/EbonJaeger/beluga/config"
	"github.com/EbonJaeger/beluga/daemon"
)

// Start handles the "start" sub-command
var Start = &cmd.CMD{
	Name:  "start",
	Alias: "s",
	Short: "Start a new Beluga daemon",
	Args:  &StartArgs{},
	Run:   StartRun,
}

// StartArgs are the arguments to the "start" sub-command
type StartArgs struct{}

// StartRun executes the command and starts the daemon
func StartRun(r *cmd.RootCMD, c *cmd.CMD) {
	// Create the server
	srv := daemon.NewServer()

	// Bind to the socket
	if err := srv.Bind(); err != nil {
		log.Errorf("Error while binding server socket '%s': %s", config.Conf.Socket, err.Error())
		return
	}

	// Start serving
	if err := srv.Start(); err != nil {
		log.Errorf("Error starting Beluga: %s", err.Error())
		return
	}
}
