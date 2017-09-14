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

	InputStream  *portmidi.Stream
	OutputStream *portmidi.Stream

	wg *sync.WaitGroup
}

func (s *Surface) CreateSurfaceFromFile(longname string, name string) {
	s.Name = longname
	s.Layout = PushLayout{}
	raw, err := ioutil.ReadFile(filepath.Join(name, "layout.json"))
	if err != nil {
		log.Fatal("Could not load layout.json")
	}
	json.Unmarshal(raw, &s.Layout)

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
	s.InputStream = is

	go s.handle()

	o, err := m.FindDevice(s.Config.MIDI.InputRegex, "output")
	if err != nil {
		log.Print(err)
	}
	os, err := m.NewOutputStream(o)
	if err != nil {
		log.Print(err)
	}
	s.OutputStream = os

	log.Print("Surface '", s.Name, "'Bound")
}

func (s *Surface) handle() {
	defer s.InputStream.Close()
	ch := s.InputStream.Listen()
	for {
		event := <-ch
		fmt.Println(event)
	}
	s.wg.Done()
}
