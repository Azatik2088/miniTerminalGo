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
			os.Exit(0)
		}
		command = strings.TrimSpace(command)
		fmt.Printf("%s: command not found\n", command)
	}
}
