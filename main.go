package main

import (
	"fmt"
	"log"
	"os"

	"github.com/zwzn/dnd/character"
	"github.com/zwzn/dnd/event"
)

func main() {
	ch, err := character.NewFile("example.md")
	if err != nil {
		log.Fatal(err)
	}

	err = event.UpdateCharacterFile(ch, "example.log")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ch.Render(os.Stdout))

}
