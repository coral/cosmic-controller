package main

import (
	"github.com/coral/cosmic-controller/pkg/surfaces"
)

func main() {
	//midi.Start()
	s := surfaces.Surface{}
	s.LoadLayout("./data/surfaces/push/layout.json")
}
