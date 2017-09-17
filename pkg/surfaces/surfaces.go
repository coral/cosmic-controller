package surfaces

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/rakyll/portmidi"

	"github.com/coral/cosmic-controller/pkg/midi"
)

type Surface struct {
	Name   string
	Config struct {
		MIDI struct {
			Input  midi.MIDIDevice
			Output midi.MIDIDevice
		}
	}
	Layout PushLayout

	Note      map[int64]Trigger
	CC        map[int64]Trigger
	PitchBend Trigger
	Triggers  map[string]Trigger

	InputMIDIStream  *portmidi.Stream
	OutputMIDIStream *portmidi.Stream

	InputChannel <-chan portmidi.Event

	Listeners []chan<- Event

	wg *sync.WaitGroup
}

type Trigger interface {
	Handle(status string, e portmidi.Event) Event
}

type Event struct {
	EventID string
	Number  int
	Type    string
	Message string
	Name    string
	Raw     portmidi.Event
}

func (s *Surface) CreateSurfaceFromFile(longname string, name string) {
	s.Name = longname
	s.Layout = PushLayout{}
	raw, err := ioutil.ReadFile(filepath.Join(name, "layout.json"))
	if err != nil {
		log.Fatal("Could not load layout.json")
	}
	json.Unmarshal(raw, &s.Layout)

	s.Note = make(map[int64]Trigger)
	s.CC = make(map[int64]Trigger)
	s.PitchBend = PushTouchable{}
	s.Triggers = make(map[string]Trigger)

	for _, element := range s.Layout.Parts.Pads {
		t := element
		s.Note[int64(element.Number)] = t
		s.Triggers[element.Name] = t
	}

	for _, element := range s.Layout.Parts.Buttons {
		t := element
		s.CC[int64(element.Number)] = t
		s.Triggers[element.Name] = t
	}

	for _, element := range s.Layout.Parts.RotaryEncoders {
		t := element
		s.Note[int64(element.Touch.Number)] = t.Touch
		s.CC[int64(element.Number)] = t
		s.Triggers[element.Name] = t
	}

	s.Note[int64(s.Layout.Parts.Slider.Touch.Number)] = s.Layout.Parts.Slider.Touch

	t := PushConfig{}
	raw, err = ioutil.ReadFile(filepath.Join(name, "config.json"))
	if err != nil {
		log.Fatal("Could not load config.json")
	}
	json.Unmarshal(raw, &t)

	s.Config.MIDI.Input = midi.MIDIDevice{
		Name:      t.MIDI.Input,
		Regex:     regexp.MustCompile(`(` + t.MIDI.Input + `)`),
		Direction: "input",
	}

	s.Config.MIDI.Output = midi.MIDIDevice{
		Name:      t.MIDI.Output,
		Regex:     regexp.MustCompile(`(` + t.MIDI.Output + `)`),
		Direction: "output",
	}

	log.Print("Surface '", s.Name, "'Loaded")
}

func (s *Surface) Bind(m midi.Handler, parentWg *sync.WaitGroup) {
	s.wg = parentWg
	i, err := m.FindDevice(s.Config.MIDI.Input)
	if err != nil {
		log.Fatal(err)
	}
	is, err := m.NewInputStream(i)
	if err != nil {
		log.Fatal(err)
	}
	s.InputMIDIStream = is
	go s.handleMessage()

	o, err := m.FindDevice(s.Config.MIDI.Output)
	if err != nil {
		log.Print(err)
	}
	os, err := m.NewOutputStream(o)
	if err != nil {
		log.Print(err)
	}
	s.OutputMIDIStream = os

	log.Print("Surface '", s.Name, "'Bound")
}

func (s *Surface) NewListener() <-chan Event {
	nc := make(chan Event)
	s.Listeners = append(s.Listeners, nc)
	return nc
}

func (s *Surface) handleMessage() {
	defer s.InputMIDIStream.Close()
	ch := s.InputMIDIStream.Listen()
	for {
		event := <-ch
		switch event.Status {
		case 144:
			//NOTE ON
			if trig, ok := s.Note[event.Data1]; ok {
				go s.processEvent(trig.Handle("On", event))
			} else {
				log.Println("Unmapped trigger:", event)
			}
		case 128:
			//NOTE OFF
			if trig, ok := s.Note[event.Data1]; ok {
				go s.processEvent(trig.Handle("Off", event))
			} else {
				log.Println("Unmapped trigger:", event)
			}
		case 176:
			//CC
			if trig, ok := s.CC[event.Data1]; ok {
				go s.processEvent(trig.Handle("CC", event))
			} else {
				log.Println("Unmapped trigger:", event)
			}
		case 224:
			//CC
			go s.processEvent(s.PitchBend.Handle("Pitchbend", event))
		default:
		}

	}
	s.wg.Done()
}

func (s *Surface) WritePushSysEx(b []byte) {
	s.OutputMIDIStream.WriteSysExBytes(portmidi.Time(), preparePushSysEx(b))
}

func (s *Surface) processEvent(e Event) {
	for _, listener := range s.Listeners {
		listener <- e
	}
}

func preparePushSysEx(b []byte) []byte {
	prefix := []byte{0xF0, 0x00, 0x21, 0x1D, 0x01, 0x01}
	suffix := []byte{0xF7}
	var sysex []byte
	sysex = append(sysex, prefix...)
	sysex = append(sysex, b...)
	sysex = append(sysex, suffix...)

	return sysex
}
