package cli

import (
	"github.com/DataDrake/cli-ng/cmd"
)

// Start handles the "start" sub-command
var Start = &cmd.CMD{
	Name:  "start",
	Alias: "s",
	Short: "Start a new Beluga client",
	Args:  &StartArgs{},
	Run:   StartRun,
}

// StartArgs are the arguments to the "start" sub-command
type StartArgs struct{}

// StartRun executes the command and starts the Client
func StartRun(r *cmd.RootCMD, c *cmd.CMD) {
	// TODO: Implement
}
