package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/vascocosta/gluon-bot-ai/cli"
	"github.com/vascocosta/gluon-bot-ai/config"
	"github.com/vascocosta/gluon-bot-ai/request"
	"github.com/vascocosta/gluon-bot-ai/utils"
)

const ConfigFile = ".gluon-bot-ai.json"

func main() {
	// Parse the CLI arguments.
	args, err := cli.ParseArgs(os.Args)
	utils.CheckError(err, utils.Exit)

	// Load and render a config from a local JSON file.
	// Render replaces placeholders with config values.
	homeDir, err := os.UserHomeDir()
	utils.CheckError(err, utils.Exit)
	cfg, err := config.LoadConfig(filepath.Join(homeDir, ConfigFile))
	utils.CheckError(err, utils.Exit)
	cfg.Render()

	// Get the API key from an environment variable.
	key, present := os.LookupEnv("GAI_KEY")
	if !present {
		log.Println("Could not find API key.")
		os.Exit(1)
	}

	// Make a request and then print the response to the std output.
	// The prompts for the request come from the CLI args and config.
	request, err := request.NewRequest(args, key, cfg)
	utils.CheckError(err, utils.Exit)
	response, err := request.Send()
	utils.CheckError(err, utils.Exit)

	for _, line := range response {
		fmt.Print(line)
	}
}
