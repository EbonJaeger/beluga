package cli

import (
	"github.com/DataDrake/cli-ng/cmd"
	"github.com/EbonJaeger/beluga/api/v1"
	"github.com/EbonJaeger/beluga/config"
)

// Stop handles the "stop" sub-command
var Stop = &cmd.CMD{
	Name:  "stop",
	Alias: "st",
	Short: "Stop a running Beluga daemon",
	Args:  &StopArgs{},
	Run:   StopRun,
}

// StopArgs are the arguments to the "stop" sub-command
type StopArgs struct{}

// StopRun will shut down the running daemon
func StopRun(r *cmd.RootCMD, c *cmd.CMD) {
	// Create an API client
	client := v1.NewClient(config.Conf.Socket)
	defer client.Close()

	// Send our request
	if err := client.Stop(); err != nil {
		log.Fatalf("Error while stopping daemon: %s\n", err.Error())
	}
	log.Goodln("Beluga has been stopped successfully")
}
