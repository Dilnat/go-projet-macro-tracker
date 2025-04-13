package main

import (
	"bufio"
	"fmt"
	"macro-tracker/cli"
	"macro-tracker/data/config"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Message struct {
	Cmd  string
	Args []string
}

func main() {
	_ = godotenv.Load()
	_ = config.ConnectDB()
	defer config.DB.Close()

	fmt.Println("Bienvenue dans Macro-Tracker CLI ✨")
	fmt.Println("💡 Tape 'help' pour voir les commandes disponibles.")

	msgChan := make(chan Message)

	// 🎧 Goroutine : lecture interactive
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("> ")
			line, _ := reader.ReadString('\n')
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			parts := strings.Fields(line)
			msgChan <- Message{Cmd: parts[0], Args: parts[1:]}
		}
	}()

	// 🔄 Traitement central
	for {
		select {
		case msg := <-msgChan:
			args := append([]string{msg.Cmd}, msg.Args...)
			cli.HandleCommand(args)
			if msg.Cmd == "exit" || msg.Cmd == "quit" {
				fmt.Println("À bientôt 👋")
				return
			}
		}
	}
}
