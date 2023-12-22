package main

import (
	"flag"

	natgridapi "github.com/jamesgawn/ng-dfs-notifier/pkg/natgridapi"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	debug := flag.Bool("debug", false, "sets the log level to debug")

	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Info().Msg("Hello World")

	natgridapi.GetDemandFlexibilityServiceRequirements()
}
