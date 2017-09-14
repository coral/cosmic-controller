package surfaces

import (
	"encoding/json"
	"fmt"
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
			InputName   string
			InputRegex  *regexp.Regexp
			OutputName  string
			OutputRegex *regexp.Regexp
		}
	}
	Layout PushLayout

	Note map[int64]Trigger
	CC   map[int64]Trigger

	InputMIDIStream  *portmidi.Stream
	OutputMIDIStream *portmidi.Stream

	InputChannel <-chan portmidi.Event

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

	for _, element := range s.Layout.Parts.Pads {
		t := element
		s.Note[int64(element.Number)] = t
	}

	for _, element := range s.Layout.Parts.Buttons {
		t := element
		s.CC[int64(element.Number)] = t
	}

	for _, element := range s.Layout.Parts.RotaryEncoders {
		t := element
		s.Note[int64(element.Touch.Number)] = t.Touch
		s.CC[int64(element.Number)] = t
	}

	s.Note[int64(s.Layout.Parts.Slider.Touch.Number)] = s.Layout.Parts.Slider.Touch

	t := PushConfig{}
	raw, err = ioutil.ReadFile(filepath.Join(name, "config.json"))
	if err != nil {
		log.Fatal("Could not load config.json")
	}
	json.Unmarshal(raw, &t)

	s.Config.MIDI.InputName = t.MIDI.Input
	s.Config.MIDI.InputRegex = regexp.MustCompile(`(` + s.Config.MIDI.InputName + `)`)
	s.Config.MIDI.OutputName = t.MIDI.Output
	s.Config.MIDI.OutputRegex = regexp.MustCompile(`(` + s.Config.MIDI.OutputName + `)`)

	log.Print("Surface '", s.Name, "'Loaded")
}

func (s *Surface) Bind(m midi.Handler, parentWg *sync.WaitGroup) {
	s.wg = parentWg
	i, err := m.FindDevice(s.Config.MIDI.InputRegex, "input")
	if err != nil {
		log.Print(err)
	}
	is, err := m.NewInputStream(i)
	if err != nil {
		log.Print(err)
	}
	s.InputMIDIStream = is
	go s.handleMessage()

	o, err := m.FindDevice(s.Config.MIDI.InputRegex, "output")
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

func (s *Surface) handleMessage() {
	defer s.InputMIDIStream.Close()
	ch := s.InputMIDIStream.Listen()
	for {
		event := <-ch

		switch event.Status {
		case 144:
			//NOTE ON
			go processEvent(s.Note[event.Data1].Handle("Note On", event))
		case 128:
			//NOTE OFF
			go processEvent(s.Note[event.Data1].Handle("Note Off", event))
		case 176:
			//CC
			switch event.Data2 {
			case 0:
				//CC ON
				go processEvent(s.CC[event.Data1].Handle("CC On", event))
			case 127:
				//CC OFF
				go processEvent(s.CC[event.Data1].Handle("CC Off", event))
			}
		default:
		}
		// if event.Status == 144 {
		// 	s.handleNote(s.Note[event.Data1], event)
		// }
		// if event.Status == 176 {
		// 	s.handleCC(s.Note[event.Data1], event)
		// }
	}
	s.wg.Done()
}

func processEvent(e Event) {
	fmt.Println(e)
}
