package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/diamondburned/arikawa/v2/state"
)

func init() {
	log.SetFlags(0)

	flag.Usage = func() {
		log.Println("usage:", filepath.Base(os.Args[0]), "-t \"TOKEN\" [-u USERID]")
		flag.PrintDefaults()
	}
}

func main() {
	token := flag.String("t", "", "discord account token")
	userid := flag.Int("u", 0, "group owner id")
	flag.Parse()

	var raider discord.UserID

	if *userid != 0 {
		raider = discord.UserID(*userid)
	}

	if *token == "" {
		log.Printf("no discord token supplied\n\n")
		log.Println("usage:", filepath.Base(os.Args[0]), "-t \"TOKEN\" [-u USERID]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	s, err := state.New(*token)
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
