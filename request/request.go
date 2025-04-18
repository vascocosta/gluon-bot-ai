package request

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/vascocosta/gluon-bot-ai/cli"
	"github.com/vascocosta/gluon-bot-ai/config"
	"github.com/vascocosta/gluon-bot-ai/utils"
	"google.golang.org/api/option"
)

type Request struct {
	cfg          *config.Config
	key          string
	systemPrompt []genai.Part
	userPrompt   genai.Text
}

func NewRequest(args cli.Args, key string, cfg *config.Config) (*Request, error) {
	systemPrompt := utils.Map(cfg.SystemPrompt, func(s string) genai.Part {
		return genai.Text(s)
	})
	chatLog, err := utils.ReadLastNLines(cfg.LogsPath+args.Channel+".txt", cfg.LogLines)
	if err != nil {
		return nil, err
	}
	userPrompt := strings.Join(append(chatLog, args.Prompt), "\n")

	return &Request{
		cfg:          cfg,
		key:          key,
		systemPrompt: systemPrompt,
		userPrompt:   genai.Text(userPrompt),
	}, nil
}

func (r *Request) Send() ([]string, error) {
	var response []string

	// Create a new API client with a root context and API key.
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(r.key))
	if err != nil {
		return response, err
	}
	defer client.Close()

	// Create a new generative model as per the config, configure it with the
	// system prompt and generate content according to the provided prompt.
	model := client.GenerativeModel(r.cfg.ModelName)
	model.SystemInstruction = genai.NewUserContent(r.systemPrompt...)
	model.SetTemperature(r.cfg.ModelTemp)
	resp, err := model.GenerateContent(ctx, r.userPrompt)
	if err != nil {
		return response, err
	}

	// Convert the generated response candidates into a slice of string.
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				response = append(response, fmt.Sprintf("%s", part))
			}
		}
	}

	return response, nil
}
