package surfaces

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Surface struct {
	layout PushLayout
}

func (s *Surface) LoadLayout(file string) {
	s.layout = PushLayout{}
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal("could not load layout")
	}
	json.Unmarshal(raw, &s.layout)
	log.Print("Layout loaded")
}
