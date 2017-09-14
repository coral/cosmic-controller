package surfaces

type PushConfig struct {
	MIDI struct {
		Input  string `json:"input"`
		Output string `json:"output"`
	} `json:"midi"`
}

type PushLayout struct {
	Parts struct {
		Pads           []PushPressable     `json:"pads"`
		Buttons        []PushPressable     `json:"buttons"`
		RotaryEncoders []PushRotaryEncoder `json:"rotaryencoders"`
		PitchBend      struct {
			Touch PushTouchable `json:"Touch"`
		} `json:"pitchbend"`
	} `json:"parts"`
	Layout struct {
		XY [][]int `json:"xy"`
	} `json:"layout"`
}

type PushPressable struct {
	Number  int    `json:"Number"`
	Type    string `json:"Type"`
	Message string `json:"Message"`
	Color   bool   `json:"Color"`
	Name    string `json:"Name"`
}

type PushRotaryEncoder struct {
	Number  int           `json:"Number"`
	Type    string        `json:"Type"`
	Message string        `json:"Message"`
	Touch   PushTouchable `json:"Touch"`
}

type PushTouchable struct {
	Number  int    `json:"Number"`
	Message string `json:"Type"`
}
