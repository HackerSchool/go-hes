package main

import (
	"errors"
	"github.com/Jguer/go-hes/driver"
	keybd "github.com/micmonay/keybd_event"
	"io/ioutil"
	"log"
	"strings"
	"sync"
)

type keybinding struct {
	Left   string `json:"left"`
	B      string `json:"b"`
	Start  string `json:"start"`
	Select string `json:"select"`
	A      string `json:"a"`
	Right  string `json:"right"`
	Up     string `json:"up"`
	Down   string `json:"down"`
}

// findArduino looks for the file that represents the Arduino
// serial connection. Returns the fully qualified path to the
// device if we are able to find a likely candidate for an
// Arduino, otherwise an empty string if unable to find
// something that 'looks' like an Arduino device.
func findArduino() ([]string, int, error) {
	contents, _ := ioutil.ReadDir("/dev")
	var n int
	var duinos []string
	// Look for what is mostly likely the Arduino device
	for _, f := range contents {
		if strings.Contains(f.Name(), "tty.usbserial") ||
			strings.Contains(f.Name(), "ttyUSB") {
			duinos = append(duinos, "/dev/"+f.Name())
			// log.Println("/dev/" + f.Name())
			n++
		}
	}

	if n != 0 {
		return duinos, n, nil
	}

	// Have not been able to find a USB device that 'looks'
	// like an Arduino.
	return duinos, n, errors.New("Device Find: Unable to find HES")
}

// translateKeybindings converts strings from keybinding struct to keybd identifiers
func translateKeybindings(kb keybinding) [8]int {
	keymap := map[string]int{
		"space": keybd.VK_ENTER,
		"a":     keybd.VK_A,
		"b":     keybd.VK_B,
		"esc":   keybd.VK_ESC,
		"c":     keybd.VK_C,
		"v":     keybd.VK_V,
		"d":     keybd.VK_D,
		"w":     keybd.VK_W,
		"s":     keybd.VK_S,
		"l":     keybd.VK_L,
		"k":     keybd.VK_K,
		"m":     keybd.VK_M,
		"n":     keybd.VK_N,
		"g":     keybd.VK_G,
		"up":    keybd.VK_UP,
		"right": keybd.VK_RIGHT,
		"down":  keybd.VK_DOWN,
		"left":  keybd.VK_LEFT,
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

func main() {
	var wg sync.WaitGroup

	kbds, err := readProfile()
	// Find the device that represents the arduino serial
	// connection.
	duinos, n, err := findArduino()
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < n; i++ {
		//log.Printf("%d %s %+v\n", i, duinos[i], kbds[i])
		wg.Add(1)
		go driver.CreateController(duinos[i], translateKeybindings(kbds[i]), &wg)
	}
	wg.Wait()
}
