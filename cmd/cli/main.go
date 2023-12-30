/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"cobra_starter/cmd/cli/cmd"
	"cobra_starter/cmd/cli/shared"
	"context"
	"os"

	"github.com/rs/zerolog"
)

func main() {

	log := zerolog.New(os.Stdout).With().Caller().Timestamp().Logger()
	ctx := context.Background()
	ctx = log.WithContext(ctx)
	log.Info().Msg("Starting application")
	shared.Version = version
	cmd.Execute(ctx)
}

var version = "Development"
