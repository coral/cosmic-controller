package surfaces

import (
	"github.com/rakyll/portmidi"
)

type PushConfig struct {
	MIDI struct {
		Input  string `json:"input"`
		Output string `json:"output"`
	} `json:"midi"`
}

type PushLayout struct {
	Parts struct {
		Pads           []PushPressable     `json:"Pads"`
		Buttons        []PushPressable     `json:"Buttons"`
		RotaryEncoders []PushRotaryEncoder `json:"RotaryEncoders"`
		Slider         struct {
			Touch PushTouchable `json:"Touch"`
		} `json:"Slider"`
	} `json:"Parts"`
	Layout struct {
		XY [][]int `json:"XY"`
	} `json:"layout"`
}

type PushPressable struct {
	Number  int    `json:"Number"`
	Type    string `json:"Type"`
	Message string `json:"Message"`
	Color   bool   `json:"Color"`
	Name    string `json:"Name"`
}

func (pt PushPressable) Handle(status string, e portmidi.Event) Event {
	st := status
	if status == "CC" {
		if e.Data2 == 127 {
			st = "On"
		} else {
			st = "Off"
		}
	}
	return Event{
		EventID: st,
		Number:  pt.Number,
		Type:    pt.Type,
		Message: pt.Message,
		Name:    pt.Name,
		Raw:     e,
	}
}

type PushRotaryEncoder struct {
	Number  int           `json:"Number"`
	Type    string        `json:"Type"`
	Message string        `json:"Message"`
	Name    string        `json:"Name"`
	Touch   PushTouchable `json:"Touch"`
}

func (pt PushRotaryEncoder) Handle(status string, e portmidi.Event) Event {
	var st string
	if e.Data2 < 64 {
		st = "increment"
	} else {
		st = "decrement"
	}
	return Event{
		EventID: st,
		Number:  pt.Number,
		Type:    pt.Type,
		Message: pt.Message,
		Name:    pt.Name,
		Raw:     e,
	}
}

type PushTouchable struct {
	Number  int    `json:"Number"`
	Message string `json:"Message"`
}

func (pt PushTouchable) Handle(status string, e portmidi.Event) Event {
	return Event{
		EventID: status,
		Number:  pt.Number,
		Type:    "note",
		Message: pt.Message,
		Name:    "Touch",
		Raw:     e,
	}
}
