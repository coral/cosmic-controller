package layers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/coral/cosmic-controller/pkg/surfaces"
)

type LayerManager struct {
	Input    <-chan surfaces.Event
	Layers   []Layer
	Surface  *surfaces.Surface
	Bindings struct {
		XY map[int]map[int]Binding
	}
}

func (lm *LayerManager) Initalize(s *surfaces.Surface) {
	lm.Surface = s
	lm.Input = s.NewListener()
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

func (lm *LayerManager) LoadBindingsFromFile(name string) {

	bf := BindingFile{}
	raw, err := ioutil.ReadFile(filepath.Join(name, "binding.json"))
	if err != nil {
		log.Fatal("Could not load layout.json")
	}
	json.Unmarshal(raw, &bf)

	//LOAD XY MAP
	newXY := make(map[int]map[int]Binding)
	for i, x := range bf.XY {
		newXY[i+1] = make(map[int]Binding)
		for j, y := range x {
			newXY[i+1][j+1] = Binding{
				Name:    y.Name,
				Trigger: lm.Surface.Triggers[y.Bind],
			}
		}
	}
	lm.Bindings.XY = newXY
	fmt.Println(newXY[2][4])
}
