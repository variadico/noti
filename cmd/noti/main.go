package main

import (
	"log"

	"github.com/variadico/noti/internal/command"
)

func main() {
	command.InitFlags(command.Root.Flags())
	if err := command.Root.Execute(); err != nil {
		log.Fatal(err)
	}
}
