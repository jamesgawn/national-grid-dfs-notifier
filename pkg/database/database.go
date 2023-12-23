package database

import (
	"database/sql"
	"os"

	"github.com/jamesgawn/ng-dfs-notifier/pkg/natgridapi"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

type DFSDatabase struct {
	db *sql.DB
}

func Open(databaseFile string) (*DFSDatabase, error) {
	wd, err := os.Getwd()
	log.Info().Str("databaseLocation", databaseFile).Str("workingDirectory", wd).Msg("Opening database file")
	db, err := sql.Open("sqlite3", databaseFile)
	if err != nil {
		log.Error().Err(err).Msg("Failed to open the database")
		return nil, err
	}

	database := &DFSDatabase{
		db: db,
	}

	err = database.initDatabaseIfNeeded()
	if err != nil {
		return nil, err
	}
	return database, nil
}

func (dfsDatabase DFSDatabase) initDatabaseIfNeeded() error {
	const create string = `
	CREATE TABLE IF NOT EXISTS requirements (
		starts TIMESTAMP NOT NULL,
		ends TIMESTAMP NOT NULL
	)`

	_, err := dfsDatabase.db.Exec(create)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create database")
		return err
	}

	return nil
}

func (dfsDatabase DFSDatabase) CheckIfRequirementHasBeenNotified(requirement natgridapi.DFSRequrement) (bool, error) {
	rows, err := dfsDatabase.db.Query("SELECT * FROM requirements WHERE starts = ? AND ends = ?", requirement.To, requirement.From)
	if err != nil {
		log.Error().Err(err).Msg("Failed to query to find the requirement")
		return false, err
	}

	defer rows.Close()

	isRow := rows.Next()
	if isRow {
		log.Info().Int("requirementId", requirement.Id).Msg("Has previously existing entry, so has previously been notified.")
	} else {
		log.Info().Int("requirementId", requirement.Id).Msg("No previously existing entry, so has not previously been notified.")
	}

	return isRow, nil
}

func (dfsDatabase DFSDatabase) AddRequirement(requirement natgridapi.DFSRequrement) error {
	_, err := dfsDatabase.db.Exec("INSERT INTO requirements VALUES(?, ?);", requirement.To, requirement.From)
	if err != nil {
		log.Error().Err(err).Int("requirementId", requirement.Id).Msg("Failed to add requirement to database")
		return err
	}
	return nil
}
