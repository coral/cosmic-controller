package layers

import "github.com/coral/cosmic-controller/pkg/surfaces"

type Binding struct {
	Name    string
	Trigger surfaces.Trigger
}

type BindingFile struct {
	XY [][]BindingEntry `json:"XY"`
}

type BindingEntry struct {
	Name string `json:"Name"`
	Bind string `json:"Bind"`
}
