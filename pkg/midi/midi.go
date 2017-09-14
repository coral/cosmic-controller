package midi

import (
	"fmt"
	"log"

	"github.com/rakyll/portmidi"
)

var (
	Device int = 2
)

func Start() {
	portmidi.Initialize()
	fmt.Println(portmidi.CountDevices())
	for i := 0; i < portmidi.CountDevices(); i++ {
		fmt.Println(i, portmidi.Info(portmidi.DeviceID(i)))
	}
	out, err := portmidi.NewOutputStream(portmidi.DeviceID(5), 1024, 0)
	if err != nil {
		log.Fatal(err)
	}
	out.WriteShort(0x90, 36, 0x01)
	out.WriteShort(0xB0, 20, 0x01)
	out.WriteShort(0xB0, 27, 0x01)
	out.WriteShort(0x90, 60, 100)

	in, err := portmidi.NewInputStream(portmidi.DeviceID(Device), 1024)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	ch := in.Listen()
	for {
		event := <-ch
		fmt.Println(event)
	}

}
