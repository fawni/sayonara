package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/diamondburned/arikawa/v2/state"
)

func main() {
	args := os.Args[1:]
	token := args[0]
	var raider discord.UserID
	if len(args) > 1 {
		s, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatalln("failed to convert userid to int:", err)
		}
		raider = discord.UserID(s)
	}
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

	groups := []discord.Channel{}
	for _, dm := range s.Ready().PrivateChannels {
		if (dm.Type != 3) || (raider != 0 && dm.DMOwnerID != raider) {
			continue
		}
		fmt.Printf("  - %s (%d)\n", dm.Name, dm.ID)
		groups = append(groups, dm)
		continue
	}

	if !ask("continue? (y|n)", 'Y', 'y') {
		fmt.Println()
		log.Fatalln("cancelled")
	}

	rand.Seed(time.Now().UnixNano())
	for _, group := range groups {
		if err := s.DeleteChannel(group.ID); err != nil {
			log.Printf("failed to leave group %d: %v", group.ID, err)
		}
		time.Sleep(time.Duration(rand.Intn(500)))
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
