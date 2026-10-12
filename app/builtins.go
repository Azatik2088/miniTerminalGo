package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func handleExit(s *Shell, args []string) {
	os.Exit(0)
	writeHistoryFile(s, s.historyFile)
}

func handleEcho(_ *Shell, args []string) {
	fmt.Println(strings.Join(args, " "))
}

func handlePwd(_ *Shell, args []string) {
	currentDir, err := os.Getwd()
	if err != nil {
		return
	}
	fmt.Println(currentDir)
}

func handleCd(_ *Shell, args []string) {
	var targetDir string
	if len(args) < 1 || args[0] == "~" {
		targetDir = os.Getenv("HOME")
	} else {
		targetDir = args[0]
	}

	if err := os.Chdir(targetDir); err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", targetDir)
	}
}

func handleType(s *Shell, args []string) {
	if len(args) == 0 {
		return
	}

	if _, ok := s.builtins[args[0]]; ok {
		fmt.Printf("%s is a shell builtin\n", args[0])
	} else if path, err := exec.LookPath(args[0]); err == nil {
		fmt.Printf("%s is %s\n", args[0], path)
	} else {
		fmt.Printf("%s: not found\n", args[0])
	}
}

func handleHistory(s *Shell, args []string) {
	if len(args) == 0 {
		printHistory(s)
		return
	}

	switch args[0] {
	case "-r":
		if len(args) < 2 {
			fmt.Println("history: -r: option requires an argument")
			return
		}
		readHistoryFile(s, args[1])
	case "-w":
		if len(args) < 2 {
			fmt.Println("history: -w: option requires an argument")
			return
		}
		writeHistoryFile(s, args[1])
	case "-a":
		if len(args) < 2 {
			fmt.Println("history: -a: option requires an argument")
			return
		}
		appendHistoryFile(s, args[1])
	}
}

func printHistory(s *Shell) {
	for index, command := range s.history {
		fmt.Printf("%5d  %s\n", index+1, command)
	}
}

func readHistoryFile(s *Shell, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()
	readHistory := bufio.NewScanner(file)

	for readHistory.Scan() {
		line := readHistory.Text()
		if strings.TrimSpace(line) != "" {
			s.history = append(s.history, line)
		}
	}
	if err := readHistory.Err(); err != nil {
		fmt.Println("Error reading history file: ", err)
		return
	}

}

func writeHistoryFile(s *Shell, filePath string) {
	content := strings.Join(s.history, "\n") + "\n"
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
}

func appendHistoryFile(s *Shell, filePath string) {
	newCommands := s.history[s.lastAppended:]
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer file.Close()

	content := strings.Join(newCommands, "\n") + "\n"
	if _, err := file.WriteString(content); err != nil {
		fmt.Println("Error: ", err)
		return
	}
	s.lastAppended = len(s.history)
}
