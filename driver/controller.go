package driver

import (
	"bufio"
	"log"
	"runtime"
	"strconv"
	"time"

	keybd "github.com/jguer/keybd_event"
	s "go.bug.st/serial.v1"
)

const repeatDelay time.Duration = 150 //Milliseconds
const sleepTime time.Duration = 1200  //Milliseconds

// CreateController reads all valid serial ports
// and handles all communication and key interpretation.
func CreateController(port s.Port, kbArray [8]int, exit chan bool) {
	// Creates Keyboard
	kb, err := keybd.NewKeyBonding()
	if err != nil {
		log.Fatal(err)
	}
	// Arduino and Keyboard Setup time
	if runtime.GOOS == "linux" {
		time.Sleep(sleepTime * time.Millisecond)
	}

	reader := bufio.NewReader(port)
	//Each HES key has its on channel to signal completion
	var signal [8]chan bool
	index := 0
	counter := 0

	for i := range signal {
		signal[i] = make(chan bool)
	}

	for {
		readKey, err := reader.ReadBytes('\n')
		if err != nil {
			panic(err)
		}

		if readKey[0] == 'P' {
			// log.Println("Pressed")
			index, _ = strconv.Atoi(string(readKey[1]))

			//Pressing Select 5 Times closes the controller
			if index == 0 {
				counter++
				if counter == 5 {
					port.Close()
					exit <- true
				}
			} else {
				counter = 0
			}
			go sendKeys(signal[index], readKey[1], kb, kbArray)
		} else if readKey[0] == 'R' {
			index, _ = strconv.Atoi(string(readKey[1]))
			signal[index] <- true
			// log.Println("Released")
		}
		// log.Println(string(read_key))
	}
}

// gkey converts the received key to a keyboard key code.
func gkey(key byte, kbArray [8]int) int {
	x, _ := strconv.Atoi(string(key))
	return kbArray[x]
}

// sendKeys handles sending keystrokes to host system.
func sendKeys(signal chan bool, key byte, kb keybd.KeyBonding, kbArray [8]int) {
	kb.SetKeys(gkey(key, kbArray)) //set keys
	var err error
	for {
		select {
		case <-signal:
			return
		default:
			// log.Println("Sent a keystroke " + string(key))
			err = kb.Launching() //launch
			if err != nil {
				panic(err)
			}
			time.Sleep(repeatDelay * time.Millisecond)
		}
	}
}
