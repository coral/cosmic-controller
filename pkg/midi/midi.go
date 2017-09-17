package midi

import (
	"fmt"
	"log"
	"regexp"

	"github.com/rakyll/portmidi"
)

type Handler struct {
}

type MIDIDevice struct {
	Name      string
	Regex     *regexp.Regexp
	Direction string
}

func Create() Handler {
	portmidi.Initialize()
	v := Handler{}
	log.Println("MIDI initialized")
	log.Println("Found these MIDI Devices:")
	t := v.GetDevices()

	for i, e := range t {

		log.Println("DEVICE:", i, e.Name)
	}
	return v
}

func (h *Handler) GetDevices() map[int]*portmidi.DeviceInfo {
	r := make(map[int]*portmidi.DeviceInfo)
	for i := 0; i < portmidi.CountDevices(); i++ {

		r[i] = portmidi.Info(portmidi.DeviceID(i))
	}
	return r
}

func (h *Handler) FindDevice(dc MIDIDevice) (portmidi.DeviceID, error) {

	for i := 0; i < portmidi.CountDevices(); i++ {
		d := portmidi.Info(portmidi.DeviceID(i))
		if dc.Regex.MatchString(d.Name) {
			if dc.Direction == "input" && d.IsInputAvailable {
				return portmidi.DeviceID(i), nil
			}
			if dc.Direction == "output" && d.IsOutputAvailable {
				return portmidi.DeviceID(i), nil
			}
		}
	}
	return portmidi.DeviceID(0), fmt.Errorf("Could not find device ")
}

func (h *Handler) NewOutputStream(id portmidi.DeviceID) (*portmidi.Stream, error) {
	s, err := portmidi.NewOutputStream(id, 1024, 0)
	if err != nil {
		return nil, err
	}
	log.Print("Opening ", portmidi.Info(id).Name, " as output")
	return s, nil
}

func (h *Handler) NewInputStream(id portmidi.DeviceID) (*portmidi.Stream, error) {
	s, err := portmidi.NewInputStream(id, 1024)
	if err != nil {
		return nil, err
	}
	log.Print("Opening ", portmidi.Info(id).Name, " as input")
	return s, nil
}

// func Start() {
// 	rand.Seed(time.Now().Unix())
// 	portmidi.Initialize()
// 	fmt.Println(portmidi.CountDevices())
// 	for i := 0; i < portmidi.CountDevices(); i++ {
// 		fmt.Println(i, portmidi.Info(portmidi.DeviceID(i)))
// 	}
// 	out, err := portmidi.NewOutputStream(portmidi.DeviceID(5), 1024, 0)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	out.WriteShort(144, 99, 127)
// 	out.WriteShort(144, 36, 125)
// 	out.WriteShort(176, 60, 125)
// 	out.WriteShort(177, 60, 125)
// 	out.WriteSysExBytes(portmidi.Time(), []byte{0xF0, 0x00, 0x21, 0x1D, 0x01, 0x01, 0x04, 120, 0xF7})
// 	in, err := portmidi.NewInputStream(portmidi.DeviceID(Device), 1024)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer in.Close()

// 	ch := in.Listen()
// 	for {
// 		event := <-ch
// 		out.WriteShort(event.Status, event.Data1, int64(random(0, 127)))

// 	}

// }

// func random(min, max int) int {
// 	return rand.Intn(max-min) + min
// }
