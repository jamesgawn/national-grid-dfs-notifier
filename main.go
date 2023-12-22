package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	natgridapi "github.com/jamesgawn/ng-dfs-notifier/pkg/natgridapi"
	telegram "github.com/jamesgawn/ng-dfs-notifier/pkg/telegram"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	debug := flag.Bool("debug", false, "sets the log level to debug")
	supplierName := flag.String("supplierName", "OCTOPUS ENERGY LIMITED", "specifies for which supplier to obtain requirements")
	telegramChatId := flag.String("telegramChatId", "-4074517448", "The ID of the chat group to send the notifications")
	telegramBotToken := flag.String("telgramBotToken", "", "the bot token to send messages if new requirements are found")

	tgt := os.Getenv("TELEGRAM_BOT_TOKEN")
	if tgt != "" {
		telegramBotToken = &tgt
	}

	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Info().Msg("Searching for supplier requirements with National Grid API")

	requirements := natgridapi.GetDFSRequirementsForSupplier(*supplierName)
	futureRequirements := make([]natgridapi.DFSRequrement, 0)
	timeNow, _ := time.Parse(time.DateOnly, "2023-12-16")

	log.Debug().Msg("Searching for requirements in the future")
	for _, requirement := range requirements {
		if requirement.To.After(timeNow) {
			log.Debug().Interface("futureRequirement", requirement).Msg("Found future requirement")
			futureRequirements = append(futureRequirements, requirement)
		}
	}

	log.Info().Interface("futureRequirements", futureRequirements).Msg("Done")

	telegram := telegram.Telegram{
		Token: *telegramBotToken,
	}

	date := futureRequirements[0].From.Format(time.DateOnly)
	from := futureRequirements[0].From.Format(time.TimeOnly)
	to := futureRequirements[0].To.Format(time.TimeOnly)
	message := fmt.Sprintf("Upcoming demand on %s %s to %s", date, from, to)

	err := telegram.SendMessage(*telegramChatId, message)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to send message to Telegram")
	}

	log.Info().Str("message", message).Msg("Sent requirement to telegram")
}
