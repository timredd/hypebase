package main

import (
	"context"
	"flag"
	"os"

	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"hypebase/ent"
	"hypebase/web"
)

func main() {
	setupLogging()
	setupConfig()
	db := setupDB(true)

	server := web.NewServer(db)
	err := server.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("error running server")
	}
}

func setupLogging() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	debug := flag.Bool("debug", false, "sets log level to debug")

	flag.Parse()

	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func setupConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal().Msgf("error loading config: %e", err)
	}
}

func setupDB(migrate bool) *ent.Client {
	pgURI := viper.GetString("POSTGRES_URI")

	client, err := ent.Open("postgres", pgURI)
	if err != nil {
		log.Fatal().Err(err).Msg("error initializing database")
	}

	if migrate {
		err = client.Schema.Create(
			context.Background(),
			schema.WithDropIndex(true),
			schema.WithDropColumn(true))
		if err != nil {
			log.Fatal().Err(err).Msg("error initializing database")
		}
	}

	return client
}
