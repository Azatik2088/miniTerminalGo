package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type CommandHandler func(s *Shell, args []string)

type Shell struct {
	history      []string
	builtins     map[string]CommandHandler
	lastAppended int
	historyFile  string
}

func NewShell() *Shell {
	s := &Shell{
		history: make([]string, 0),
		builtins: map[string]CommandHandler{
			"exit":    handleExit,
			"echo":    handleEcho,
			"pwd":     handlePwd,
			"cd":      handleCd,
			"type":    handleType,
			"history": handleHistory,
		},
		lastAppended: 0,
		historyFile:  os.Getenv("HISTFILE"),
	}
	if s.historyFile != "" {
		readHistoryFile(s, s.historyFile)
		s.lastAppended = len(s.history)
	}
	return s
}

func (s *Shell) handlerExternal(cmd string, args []string) {
	if _, err := exec.LookPath(cmd); err != nil {
		fmt.Printf("%s: command not found\n", cmd)
		return
	}
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Run()
}

func (s *Shell) Run() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		command, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		command = strings.TrimSpace(command)
		parts := parseCommand(command)
		if len(parts) == 0 {
			continue
		}
		s.history = append(s.history, command)
		commandName := parts[0]
		args := parts[1:]

		if handler, ok := s.builtins[commandName]; ok {
			handler(s, args)
		} else {
			s.handlerExternal(commandName, args)
		}
	}
}

func parseCommand(command string) []string {
	result := []string{}
	var currentWord strings.Builder
	isSingleQuotes := false
	isDoubleQuotes := false
	for _, ch := range command {
		if ch == '\'' && !isDoubleQuotes {
			isSingleQuotes = !isSingleQuotes
			continue
		} else if ch == '"' && !isSingleQuotes {
			isDoubleQuotes = !isDoubleQuotes
			continue
		} else if !isSingleQuotes && !isDoubleQuotes && ch == ' ' {
			if currentWord.Len() > 0 {
				result = append(result, currentWord.String())
				currentWord.Reset()
			}
		} else {
			currentWord.WriteRune(ch)
		}
	}
	if currentWord.Len() > 0 {
		result = append(result, currentWord.String())
	}
	return result
}
