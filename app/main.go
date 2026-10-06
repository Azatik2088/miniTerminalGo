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
	shellCommands := []string{"exit", "echo", "type"}
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

		default:
			fmt.Printf("%s: command not found\n", parts[0])
		}
	}
}
