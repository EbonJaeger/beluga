package main

import (
	"os"

	"github.com/EbonJaeger/beluga"
	"github.com/jessevdk/go-flags"
)

func main() {
	var opts beluga.Flags
	parser := flags.NewParser(&opts, flags.Default)
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
	beluga.NewBeluga(opts)
}
