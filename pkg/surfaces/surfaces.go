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
			Input  midi.MIDIDevice
			Output midi.MIDIDevice
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

func (s *Surface) handleMessage() {
	defer s.InputMIDIStream.Close()
	ch := s.InputMIDIStream.Listen()
	for {
		event := <-ch

		switch event.Status {
		case 144:
			//NOTE ON
			go s.processEvent(s.Note[event.Data1].Handle("On", event))
		case 128:
			//NOTE OFF
			go s.processEvent(s.Note[event.Data1].Handle("Off", event))
		case 176:
			//CC
			go s.processEvent(s.CC[event.Data1].Handle("CC", event))
		default:
		}

	}
	s.wg.Done()
}

func (s *Surface) WritePushSysEx(b []byte) {
	s.OutputMIDIStream.WriteSysExBytes(portmidi.Time(), preparePushSysEx(b))
}

func (s *Surface) processEvent(e Event) {
	fmt.Println(e)
	//s.OutputMIDIStream.WriteShort(144, int64(e.Number), 125)

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
