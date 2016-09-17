package main

import (
	"errors"
	"github.com/Jguer/go-hes/driver"
	keybd "github.com/Jguer/keybd_event"
	"go.bug.st/serial"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"
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

const challenge string = "Hi. Who are you?"

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
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}

	mode := &serial.Mode{
		BaudRate: 9600,
	}

	var i int
	resP := make(chan *serial.SerialPort, 1)
	for _, portName := range ports {
		log.Println("Attempting connection to " + portName)
		go func() {
			port, err := serial.OpenPort(portName, mode)
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

			if strings.Contains(string(buff[:n]), "Hi. I'm HES") {
				resP <- port
			}
			resP <- nil
		}()

		select {
		case port := <-resP:
			if port != nil {
				log.Printf("Communication established with %v\n", portName)
				wg.Add(1)
				go driver.CreateController(port, translateKeybindings(kbds[i]), &wg)
				i++
			}
		case <-time.After(time.Second * 5):
			log.Println("Connection timed out on " + portName)
			continue
		}
	}
	wg.Wait()
}
