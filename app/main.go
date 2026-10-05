package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
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
			result := strings.Join(parts[1:], " ")
			fmt.Println(result)
		default:
			fmt.Printf("%s: command not found\n", command)
		}
	}
}
