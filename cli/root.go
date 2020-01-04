package cli

import (
	"github.com/DataDrake/cli-ng/cmd"
	"github.com/DataDrake/waterlog"
	"github.com/DataDrake/waterlog/format"
	"github.com/DataDrake/waterlog/level"
	log2 "log"
	"os"
)

var log *waterlog.WaterLog

// GlobalFlags holds the flags for all commands
type GlobalFlags struct{}

// Root is the main command for the application
var Root *cmd.RootCMD

func init() {
	Root = &cmd.RootCMD{
		Name:  "beluga",
		Short: "Fun Discord bot",
		Flags: &GlobalFlags{},
	}

	// Initialize subcommands
	Root.RegisterCMD(Start)

	// Initialize logging
	log = waterlog.New(os.Stdout, "", log2.Ltime)
	log.SetLevel(level.Info)
	log.SetFormat(format.Min)
}
