package main

import (
	"sync"

	"github.com/coral/cosmic-controller/pkg/midi"
	"github.com/coral/cosmic-controller/pkg/surfaces"
)

func main() {

	var wg sync.WaitGroup

	m := midi.Create()

	s := surfaces.Surface{}
	s.CreateSurfaceFromFile("Ableton Push 2", "./data/surfaces/push")
	wg.Add(1)
	s.Bind(m, &wg)

	wg.Wait()
}
