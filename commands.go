package marmot

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type commandsHolder struct {
	commands map[state][]string
	currentCommandBuilder strings.Builder
}

func getCommandsByFile(filePath string) map[state][]string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error when try read file %s: %s", filePath, err)
	}

	holder := &commandsHolder{
		commands: make(map[state][]string, 2),
		currentCommandBuilder: strings.Builder{},
	}
	var currentState state
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "-- +") {
			if strings.HasPrefix(line, "-- +Up") {
				currentState = Up
			} else if strings.HasPrefix(line, "-- +Down") {
				holder.keep(Up)
				currentState = Down
			} else if strings.HasPrefix(line, "-- +Then") {
				holder.keep(currentState)
			} else if strings.HasPrefix(line, "-- +End") {
				holder.keep(currentState)
			} else {
				log.Fatalf("Unknown migration state line: %s", line)
			}
			continue
		}
		if currentState != Up && currentState != Down {
			continue
		}
		holder.appendSql(line)
	}
	return holder.commands
}

func (c *commandsHolder) keep(state state) {
	command := c.currentCommandBuilder.String()
	if len(command) > 0 {
		c.commands[state] = append(c.commands[state], command)
	}
	c.currentCommandBuilder.Reset()
}

func (c *commandsHolder) appendSql(sqlPart string) {
	c.currentCommandBuilder.WriteString(sqlPart)
}