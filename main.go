package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	app := &cli.App{
		Name:    "tig",
		Usage:   "stupid content tracker",
		Version: "v0.0.0",
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatal(err)
	}
}
