package main

import (
	"fmt"
	"log"
	"os"

	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/diamondburned/arikawa/v2/state"
)

func main() {
	token := os.Args[1]
	if token == "" {
		log.Fatalln("missing $TOKEN")
	}

	s, err := state.New(token)
	if err != nil {
		log.Fatalln("failed to create state:", err)
	}

	ready, cancel := s.ChanFor(func(v interface{}) bool {
		_, ok := v.(*gateway.ReadyEvent)
		return ok
	})
	defer cancel()

	if err := s.Open(); err != nil {
		log.Fatalln("failed to open:", err)
	}

	defer s.CloseGracefully()

	<-ready
	cancel()

	fmt.Println("found these groups:")

	for _, dm := range s.Ready().PrivateChannels {
		if dm.Type != 3 {
			continue
		}
		if err != nil {
			fmt.Printf("  - %d (error: %s)\n", dm.ID, err)
			continue
		}
		fmt.Printf("  - %s (%d)\n", dm.Name, dm.ID)
	}

	if !ask("continue?", 'Y', 'y') {
		fmt.Println()
		log.Fatalln("cancelled")
	}

	for _, dm := range s.Ready().PrivateChannels {
		if dm.Type != 3 {
			continue
		}
		if err := s.DeleteChannel(dm.ID); err != nil {
			log.Printf("failed to leave group %d: %v", dm.ID, err)
		}
	}
}

func ask(prompt string, expect ...byte) bool {
	fmt.Print(prompt, " ")

	var c [1]byte

	_, err := os.Stdin.Read(c[:])
	if err != nil {
		return false
	}

	for _, b := range expect {
		if b == c[0] {
			return true
		}
	}

	return false
}
