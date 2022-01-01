package commands

import (
	"bufio"
	"os"
	"strings"
)

type state string

const (
	Up   state = "UP"
	Down state = "DOWN"
)

type commandsBuilder struct {
	commands map[state][]string
	sb       strings.Builder
	state    state
}

func GetCommandsByFile(migrationPath string, st state) (commands []string, error error) {
	file, err := os.Open(migrationPath)
	if err != nil {
		return nil, err
	}

	builder := &commandsBuilder{
		commands: make(map[state][]string, 2),
		sb:       strings.Builder{},
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		switch line {
		case "-- +Up":
			builder.state = Up
		case "-- +Down":
			builder.endSQLCommand()
			builder.state = Down
			continue
		case "-- +Then":
			builder.endSQLCommand()
		case "-- +End":
			builder.endSQLCommand()
			return builder.commands[st], nil
		default:
			if builder.state == st {
				builder.beginSQL(line)
			}
		}
	}
	return []string{}, nil
}

func (c *commandsBuilder) endSQLCommand() {
	command := c.sb.String()
	if len(command) > 0 {
		c.commands[c.state] = append(c.commands[c.state], command)
	}
	c.sb.Reset()
}

func (c *commandsBuilder) beginSQL(sqlPart string) {
	c.sb.WriteString(sqlPart)
}
