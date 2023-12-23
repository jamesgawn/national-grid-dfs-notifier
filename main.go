package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	database "github.com/jamesgawn/ng-dfs-notifier/pkg/database"
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
	databaseLocation := flag.String("databaseLocation", "./database.db", "The location to place the sqlite DB for the command")

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

	for _, futureRequirement := range futureRequirements {
		date := futureRequirement.From.Format(time.DateOnly)
		from := futureRequirement.From.Format(time.TimeOnly)
		to := futureRequirement.To.Format(time.TimeOnly)
		message := fmt.Sprintf("Upcoming demand on %s %s to %s", date, from, to)

		db, err := database.Open(*databaseLocation)
		if err != nil {
			log.Fatal().Err(err).Msg("Unable to open database to check if prior notification sent.")
		}
		previouslyNotified, err := db.CheckIfRequirementHasBeenNotified(futureRequirement)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to check requirements")
		}

		if previouslyNotified {
			log.Info().Int("requirementId", futureRequirement.Id).Msg("Not proceeding with notification as previously sent")
			continue
		}

		err = telegram.SendMessage(*telegramChatId, message)
		if err != nil {
			log.Fatal().Err(err).Int("requirementId", futureRequirement.Id).Msg("Failed to send message to Telegram")
		}

		log.Info().Str("message", message).Int("requirementId", futureRequirement.Id).Msg("Sent requirement to telegram")

		db.AddRequirement(futureRequirement)
		log.Info().Str("message", message).Int("requirementId", futureRequirement.Id).Msg("Logged notification in database")
	}

	log.Info().Msg("Finished processing future requirements")
}
