package main

import (
	"encoding/json"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/jguer/go-hes/driver"
	keybd "github.com/jguer/keybd_event"
	s "go.bug.st/serial.v1"
)

type keybinding struct {
	A      string `json:"a"`
	B      string `json:"b"`
	Start  string `json:"start"`
	Select string `json:"select"`
	Left   string `json:"left"`
	Right  string `json:"right"`
	Up     string `json:"up"`
	Down   string `json:"down"`
}

type serialport struct {
	name string
	sp   s.Port
}

const challenge string = "Hi. Who are you?"
const response string = "Hi. I'm HES"
const filename string = "mappings.json"

// translateKeybindings converts strings from keybinding struct to keybd identifiers
func translateKeybindings(kb keybinding) [8]int {
	keymap := map[string]int{
		"a":     keybd.VK_A,
		"b":     keybd.VK_B,
		"c":     keybd.VK_C,
		"d":     keybd.VK_D,
		"e":     keybd.VK_E,
		"f":     keybd.VK_F,
		"g":     keybd.VK_G,
		"h":     keybd.VK_H,
		"i":     keybd.VK_I,
		"j":     keybd.VK_J,
		"k":     keybd.VK_K,
		"l":     keybd.VK_L,
		"m":     keybd.VK_M,
		"n":     keybd.VK_N,
		"o":     keybd.VK_O,
		"p":     keybd.VK_P,
		"q":     keybd.VK_Q,
		"r":     keybd.VK_R,
		"s":     keybd.VK_S,
		"t":     keybd.VK_T,
		"u":     keybd.VK_U,
		"v":     keybd.VK_V,
		"w":     keybd.VK_W,
		"x":     keybd.VK_X,
		"y":     keybd.VK_Y,
		"z":     keybd.VK_Z,
		"0":     keybd.VK_0,
		"1":     keybd.VK_1,
		"2":     keybd.VK_2,
		"3":     keybd.VK_3,
		"4":     keybd.VK_4,
		"5":     keybd.VK_5,
		"6":     keybd.VK_6,
		"7":     keybd.VK_7,
		"8":     keybd.VK_8,
		"9":     keybd.VK_9,
		"up":    keybd.VK_UP,
		"down":  keybd.VK_DOWN,
		"left":  keybd.VK_LEFT,
		"right": keybd.VK_RIGHT,
		"esc":   keybd.VK_ESC,
		"space": keybd.VK_SPACE,
		"enter": keybd.VK_ENTER,
	}

	var kbArray [8]int
	kbArray[0] = keymap[kb.Select]
	kbArray[1] = keymap[kb.Start]
	kbArray[2] = keymap[kb.Up]
	kbArray[3] = keymap[kb.Down]
	kbArray[4] = keymap[kb.Left]
	kbArray[5] = keymap[kb.Right]
	kbArray[6] = keymap[kb.B]
	kbArray[7] = keymap[kb.A]

	return kbArray
}

func defaultProfile() (kbds []keybinding, err error) {
	byt := []byte(`[{"a":"v","b":"c","start":"esc","select":"b","left":"a","right":"d","up":"w","down":"s"},
        {"a":"n","b":"m","start":"l","select":"k","left":"left","right":"right","up":"up","down":"down"},
        {"a":"n","b":"m","start":"l","select":"k","left":"left","right":"right","up":"up","down":"down"},
        {"a":"n","b":"m","start":"l","select":"k","left":"left","right":"right","up":"up","down":"down"}]`)

	err = json.Unmarshal(byt, &kbds)
	if err != nil {
		return
	}

	return
}

func main() {
	probe := true
	for _, arg := range os.Args[1:] {
		if strings.Contains(arg, "config") {
			// If args include config then you'll be taken to configure the controller
			startConfig()
		} else if strings.Contains(arg, "skip") {
			// If args include skip, then the driver won't search for new controllers.
			probe = false
		}
	}

	// read keybindings
	kbds, err := readProfile()
	if err != nil {
		log.Println(err)
		kbds, err = defaultProfile()
		if err != nil {
			log.Fatal(err)
		}
	}

	mode := &s.Mode{
		BaudRate: 9600,
	}

	resP := make(chan serialport, 1)
	exit := make(chan bool, 1)
	var connections []string
	var i int

	for {
		if len(connections) == 0 || probe == true {
			// Find the device that represents the arduino serial
			// connection.
			ports, err := s.GetPortsList()
			if err != nil {
				log.Fatal(err)
			}

			for _, portName := range ports {
				// Skip unnecessary ports in linux and OSX
				if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
					if !strings.Contains(portName, "usb") && !strings.Contains(portName, "USB") {
						continue
					}
				}

				// Check if controller is already in use.
				found := false
				for _, controller := range connections {
					if controller == portName {
						found = true
					}
				}

				if found {
					continue
				}

				connections = append(connections, portName)
				go handshake(portName, mode, resP)
			}
		}

		select {
		case port := <-resP:
			if port.sp != nil {
				log.Printf("Communication established with %v\n", port.name)
				go driver.CreateController(port.sp, translateKeybindings(kbds[i]), exit)
				i++
			}
		case <-time.After(time.Second * 10):
			continue
		case <-exit:
			os.Exit(0)
		}
	}
}

func handshake(portName string, mode *s.Mode, resP chan serialport) {
	port, err := s.Open(portName, mode)
	if err != nil {
		log.Println(err)
	}

	log.Println("Executing hand shake")
	time.Sleep(1500 * time.Millisecond)
	_, err = port.Write([]byte(challenge))
	if err != nil {
		log.Println(err)
	}

	buff := make([]byte, 30)
	n, err := port.Read(buff)
	if err != nil {
		log.Println(err)
	}

	if strings.Contains(string(buff[:n]), response) {
		resP <- serialport{
			name: portName,
			sp:   port,
		}
	}

	resP <- serialport{
		name: portName,
		sp:   nil,
	}
}
