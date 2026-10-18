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

		//проверка на наличие перенаправления вывода
		terminal := os.Stdout
		var file *os.File
		parts, outputFile := findRedirection(parts)
		if outputFile != "" {
			var err error
			file, err = os.Create(outputFile)
			if err != nil {
				fmt.Println(err)
				continue
			}
			os.Stdout = file
		}

		s.history = append(s.history, command)
		commandName := parts[0]
		args := parts[1:]

		if handler, ok := s.builtins[commandName]; ok {
			handler(s, args)
		} else {
			s.handlerExternal(commandName, args)
		}
		if file != nil {
			file.Close()
			os.Stdout = terminal
		}
	}
}

func findRedirection(parts []string) (cleanParts []string, outputFile string) {
	cleanParts = parts
	for index := range parts {
		if parts[index] == ">" || parts[index] == "1>" {
			cleanParts = parts[:index]
			if index+1 < len(parts) {
				outputFile = parts[index+1]
			}
			break
		}
	}
	return
}
