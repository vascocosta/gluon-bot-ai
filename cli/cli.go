package cli

import (
	"errors"
	"strings"
)

type Args struct {
	Nick    string
	Channel string
	Prompt  string
}

func ParseArgs(args []string) (Args, error) {
	if len(args) < 4 {
		return Args{}, errors.New("Not enough parameters.")
	}

	return Args{args[1], args[2], strings.Join(args[3:], " ")}, nil
}
