package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	for {
		fmt.Print("$ ")
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			os.Exit(0)
		}
		command = strings.TrimSpace(command)
		fmt.Printf("%s: command not found\n", command)
	}
}
