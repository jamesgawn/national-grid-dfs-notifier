package main

import (
	"flag"

	ngapi "github.com/jamesgawn/ng-dfs-notifier/pkg/ng-api"
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

	ngapi.GetDemandFlexibilityServiceRequests()
}
