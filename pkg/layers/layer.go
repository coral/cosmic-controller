package layers

import "github.com/coral/cosmic-controller/pkg/surfaces"

type Layer struct {
	Name     string
	Bindings []surfaces.Trigger
}
