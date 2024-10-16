package main

import (
	"log/slog"
	"os"
	"sync"

	"gopkg.in/ini.v1"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var GITHUB_API_TOKEN string
var L1_KEYWORDS []string
var L2_KEYWORDS []string
var POSITION_KEYWORDS []string
var POSITION_ANTI_KEYWORDS []string
var LOCATION_KEYWORDS []string

func connectDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("unearth.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	migratePull(db)
	migrateUnearth(db)
	return db
}

func pullOrgsRoutine(db *gorm.DB, wg *sync.WaitGroup) {
	err := pullOrgs(GITHUB_API_TOKEN, db, nil)
	if err != nil {
		slog.Error(err.Error())
	}
	wg.Done()
}

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	pull := false
	process := false

	for _, flag := range os.Args {
		if flag == "pull" {
			pull = true
		} else if flag == "process" {
			process = true
		}
	}

	cfg, err := ini.Load("config.ini")
	if err != nil {
		slog.Error("Failed to read config file")
		os.Exit(1)
	}

	GITHUB_API_TOKEN = cfg.Section("").Key("GITHUB_API_TOKEN").String()
	L1_KEYWORDS = cfg.Section("").Key("L1LINK_KEYWORDS").Strings(",")
	L2_KEYWORDS = cfg.Section("").Key("L2LINK_KEYWORDS").Strings(",")
	POSITION_KEYWORDS = cfg.Section("").Key("POSITION_KEYWORDS").Strings(",")
	POSITION_ANTI_KEYWORDS = cfg.Section("").Key("POSITION_ANTI_KEYWORDS").Strings(",")
	LOCATION_KEYWORDS = cfg.Section("").Key("LOCATION_KEYWORDS").Strings(",")

	db := connectDatabase()

	var wg sync.WaitGroup
	if pull {
		go pullOrgsRoutine(db, &wg)
		wg.Add(1)
	}

	if process {
		unearth(db, "/tmp/unearth")
	}

	wg.Wait()
}
