package event

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/zwzn/dnd/blade"
	"github.com/zwzn/dnd/character"
)

type Event struct {
	Type string `json:"type"`
	DamageEvent
	StatusEvent
	UseEvent
	EventEvent
}
type DamageEvent struct {
	Damage int `json:"damage"`
}
type StatusEvent struct {
	Effect string `json:"effect"`
	Reset  string `json:"reset"`
}
type UseEvent struct {
	Name string `json:"name"`
}
type EventEvent struct {
	Event string `json:"event"`
}
type recharge struct {
	current int
	use     int
	total   int
}
type chWrapper struct {
	*character.Character
	status   map[string][]string
	recharge map[string]map[string]recharge
}

func UpdateCharacterFile(ch *character.Character, file string) error {

	b, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	lines := bytes.Split(b, []byte("\n"))
	events := make([]*Event, len(lines))
	for i, line := range lines {
		event := &Event{}
		json.Unmarshal(line, event)
		events[i] = event
	}
	return UpdateCharacter(ch, events)
}
func UpdateCharacter(ch *character.Character, events []*Event) error {
	chw := &chWrapper{
		Character: ch,
		status:    map[string][]string{},
		recharge:  map[string]map[string]recharge{},
	}
	chw.updateBlade()
	for _, event := range events {
		event.DamageEvent.Run(chw)
		event.StatusEvent.Run(chw)
		event.EventEvent.Run(chw)
		event.UseEvent.Run(chw)
	}
	spew.Dump(chw)
	os.Exit(1)
	return nil
}

func (c *chWrapper) updateBlade() {
	b := blade.New()
	b.Directive("recharge", func(args []string) string {
		name := args[0]
		event := args[1]
		count := 1
		use := 1
		if len(args) > 2 {
			parts := strings.Split(args[2], "/")
			if len(parts) == 1 {
				count, _ = strconv.Atoi(parts[0])
			} else {
				count, _ = strconv.Atoi(parts[0])
				use, _ = strconv.Atoi(parts[1])
				count *= use

			}
		}
		e, ok := c.recharge[event]
		if !ok {
			e = map[string]recharge{}
		}
		e[name] = recharge{
			total:   count,
			use:     use,
			current: count,
		}
		c.recharge[event] = e
		return ""
	})
	c.Blade(b)
}

func (e *DamageEvent) Run(ch *chWrapper) {
	ch.CurrentHP -= e.Damage
}
func (e *StatusEvent) Run(ch *chWrapper) {
	if e.Reset == "" || e.Effect == "" {
		return
	}
	list, ok := ch.status[e.Reset]
	if !ok {
		list = []string{}
	}
	ch.status[e.Reset] = append(list, e.Effect)
}

func (e *UseEvent) Run(ch *chWrapper) {
	if e.Name == "" {
		return
	}

	for _, recharge := range ch.recharge {
		for name, count := range recharge {
			if name != e.Name {
				continue
			}
			count.current -= count.use
			if count.current < 0 {
				count.current = 0
			}
			recharge[name] = count
		}
	}
}

func (e *EventEvent) Run(ch *chWrapper) {
	if e.Event == "" {
		return
	}
	delete(ch.status, e.Event)

	if recharge, ok := ch.recharge[e.Event]; ok {
		for name, count := range recharge {
			count.current++
			if count.current > count.total {
				count.current = count.total
			}
			recharge[name] = count
		}
	}
}
