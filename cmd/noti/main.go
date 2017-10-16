package main

import (
	"log"

	"github.com/variadico/noti/internal/command"
)

func main() {
	if err := command.Root.Execute(); err != nil {
		log.Fatal(err)
	}
}
