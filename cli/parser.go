package cli

import (
	"strings"
)

type ParsedCommand struct {
	Name string
	Args []string
}

func ParseCommand(line string) ParsedCommand {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return ParsedCommand{}
	}

	return ParsedCommand{
		Name: fields[0],
		Args: fields[1:],
	}
}
