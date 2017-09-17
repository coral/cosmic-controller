package layers

import (
	"fmt"

	"github.com/coral/cosmic-controller/pkg/surfaces"
)

type LayerManager struct {
	Input  <-chan surfaces.Event
	Layers []Layer
}

func (lm *LayerManager) Initalize(events <-chan surfaces.Event) {
	lm.Input = events
	go lm.Route()
}

func (lm *LayerManager) Route() {
	for {
		event := <-lm.Input
		fmt.Println(event)
	}
}

func (lm *LayerManager) LoadLayer(l Layer) {
	lm.Layers = append(lm.Layers, l)
}
