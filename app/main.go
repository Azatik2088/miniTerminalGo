package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	shellCommands := []string{"exit", "echo", "type", "pwd", "cd", "history"}
	var history []string
	for {
		fmt.Print("$ ")
		command, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		command = strings.TrimSpace(command)
		parts := strings.Fields(command)
		if len(parts) == 0 {
			continue
		}
		history = append(history, command)
		switch parts[0] {
		case "exit":
			return
		case "echo":
			fmt.Println(strings.Join(parts[1:], " "))
		case "type":
			if len(parts) < 2 {
				break
			}

			if slices.Contains(shellCommands, parts[1]) {
				fmt.Printf("%s is a shell builtin\n", parts[1])
			} else if path, err := exec.LookPath(parts[1]); err == nil {
				fmt.Printf("%s is %s\n", parts[1], path)
			} else {
				fmt.Printf("%s: not found\n", parts[1])
			}
		case "pwd":
			dir, err := os.Getwd()
			if err != nil {
				fmt.Println("Error: ", err)
				break
			}
			fmt.Println(dir)
		case "cd":
			var targetDir string
			if len(parts) < 2 || parts[1] == "~" {
				targetDir = os.Getenv("HOME")
			} else {
				targetDir = parts[1]
			}

			if err := os.Chdir(targetDir); err != nil {
				fmt.Printf("cd: %s: No such file or directory\n", targetDir)
			}
		case "history":
			if len(parts) > 2 && parts[1] == "-r" {
				file, err := os.Open(parts[2])
				if err != nil {
					return
				}
				defer file.Close()

				readHistory := bufio.NewScanner(file)
				for readHistory.Scan() {
					line := readHistory.Text()
					if strings.TrimSpace(line) != "" {
						history = append(history, line)
					}
				}

			} else {
				for index, command := range history {
					fmt.Printf("%5d %s\n", index+1, command)
				}
			}

		default:
			if _, err := exec.LookPath(parts[0]); err != nil {
				fmt.Printf("%s: command not found\n", parts[0])
			} else {
				cmd := exec.Command(parts[0], parts[1:]...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()
			}
		}
	}
}
