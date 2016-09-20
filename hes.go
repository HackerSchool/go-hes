package main

import (
	"github.com/Jguer/go-hes/driver"
	keybd "github.com/Jguer/keybd_event"
	"go.bug.st/serial"
	"log"
	"os"
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
		"space": keybd.VK_ENTER,
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

	for _, arg := range os.Args {
		if strings.Contains(arg, "config") {
			startConfig()
		}
	}

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
		if strings.Contains(portName, "tty") {
			continue
		}
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
		case <-time.After(time.Second * 4):
			log.Println("Connection timed out on " + portName)
			continue
		}
	}
	wg.Wait()
}
