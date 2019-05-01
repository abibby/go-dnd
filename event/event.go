package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/zwzn/go-dnd/character"
	"golang.org/x/xerrors"
)

type EventType string

const (
	LongRest  = EventType("long-rest")
	ShortRest = EventType("short-rest")
)

type Event struct {
	Type string `json:"type"`
	DamageEvent
	StatusEvent
	UseEvent
	EventEvent
	MoneyEvent
}
type DamageEvent struct {
	Damage int `json:"damage,omitempty"`
}
type StatusEvent struct {
	Effect string    `json:"effect,omitempty"`
	Reset  EventType `json:"reset,omitempty"`
}
type UseEvent struct {
	Use string `json:"use,omitempty"`
}
type EventEvent struct {
	Event EventType `json:"event,omitempty"`
}
type MoneyEvent struct {
	Amount Money `json:"amount,omitempty"`
}
type recharge struct {
	event   EventType
	current int
	use     int
	total   int
}

type Money int64

const (
	Copper   = Money(1)
	Silver   = Copper * 10
	Gold     = Silver * 10
	Platinum = Gold * 10
)

func (m Money) String() string {
	parts := []string{}
	current := m
	pp := current / Platinum
	if pp != 0 {
		parts = append(parts, fmt.Sprintf("%dpp", pp))
		current -= pp * Platinum
	}
	gp := current / Gold
	if gp != 0 {
		parts = append(parts, fmt.Sprintf("%dgp", gp))
		current -= gp * Gold
	}
	sp := current / Silver
	if sp != 0 {
		parts = append(parts, fmt.Sprintf("%dsp", sp))
		current -= sp * Silver
	}
	if current != 0 {
		parts = append(parts, fmt.Sprintf("%dcp", current))
	}
	return strings.Join(parts, " ")
}

type chWrapper struct {
	*character.Character
	status   map[EventType][]string
	recharge map[string]recharge
	cp       Money
}

func UpdateCharacterFile(ch *character.Character, file string) error {

	f, err := os.OpenFile(file, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return xerrors.Errorf("error opening log file: %w", err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return xerrors.Errorf("error reading log file: %w", err)
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
		status:    map[EventType][]string{},
		recharge:  map[string]recharge{},
	}
	chw.updatePreBlade()
	for _, event := range events {
		event.DamageEvent.Run(chw)
		event.StatusEvent.Run(chw)
		event.EventEvent.Run(chw)
		event.UseEvent.Run(chw)
		event.MoneyEvent.Run(chw)
	}
	chw.updatePostBlade()
	return nil
}

func (e *DamageEvent) Run(ch *chWrapper) {
	ch.CurrentHP -= e.Damage
}
func (e *MoneyEvent) Run(ch *chWrapper) {
	ch.cp += e.Amount
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
	if e.Use == "" {
		return
	}

	r, ok := ch.recharge[e.Use]
	if !ok {
		return
	}
	r.current -= r.use
	if r.current < 0 {
		r.current = 0
	}
	ch.recharge[e.Use] = r
}

func (e *EventEvent) Run(ch *chWrapper) {
	if e.Event == "" {
		return
	}
	delete(ch.status, e.Event)

	for name, r := range ch.recharge {
		if r.event != e.Event {
			continue
		}
		r.current += r.total / r.use
		if r.current > r.total {
			r.current = r.total
		}
		ch.recharge[name] = r
	}

	switch e.Event {
	case LongRest:
		ch.CurrentHP = ch.MaxHP
		fallthrough
	case ShortRest:

	}

}
