package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/vascocosta/gluon-bot-ai/cli"
	"github.com/vascocosta/gluon-bot-ai/config"
	"github.com/vascocosta/gluon-bot-ai/request"
	"github.com/vascocosta/gluon-bot-ai/utils"
)

const ConfigFile = ".gluon-bot-ai.json"

func main() {
	// Parse the CLI arguments.
	args, err := cli.ParseArgs(os.Args)
	utils.Try(err, utils.Exit)

	// Load and render a config from a local JSON file.
	// Render replaces placeholders with config values.
	homeDir, err := os.UserHomeDir()
	utils.Try(err, utils.Exit)
	cfg, err := config.LoadConfig(filepath.Join(homeDir, ConfigFile))
	utils.Try(err, utils.Exit)
	cfg.Render()

	// Get the API key from an environment variable.
	key, present := os.LookupEnv("GAI_KEY")
	if !present {
		log.Println("Could not get API key.")
		os.Exit(1)
	}

	// Main loop where we make a request, print a response and then wait.
	// The output file is cfg.OutPath and the interval is cfg.SleepTime.
	for {
		request, err := request.NewRequest(args, key, cfg)
		if err != nil {
			continue
		}
		response, err := request.Send()
		if err != nil {
			continue
		}

		if len(response) > 0 {
			f, err := os.Create(cfg.OutPath)
			if err != nil {
				f.Close()
				continue
			}

			_, err = fmt.Fprintf(f, "%s %s", args.Channel, response[0])
			if err != nil {
				f.Close()
				continue
			}

			f.Close()
		}

		time.Sleep(cfg.SleepTime * time.Minute)

	}
}
