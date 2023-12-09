package main

import (
	"flag"
	"net/http"

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

	response, err := http.Get("https://api.nationalgrideso.com/api/3/action/datastore_search_sql?sql=SELECT%20*%20FROM%20%20%227914dd99-fe1c-41ba-9989-5784531c58bb%22%20ORDER%20BY%20%22_id%22%20ASC%20LIMIT%20100")

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to retrieve DFS data from National Grid API")
	}

	log.Debug().Int("statusCode", response.StatusCode).Msg("Retrieved DFS data from National Grid API")
}
