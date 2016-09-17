package driver

import (
	"bufio"
	keybd "github.com/micmonay/keybd_event"
	"go.bug.st/serial"
	"log"
	"strconv"
	"sync"
	"time"
)

const repeatDelay time.Duration = 160 //Milliseconds
const sleepTime time.Duration = 1000  //Milliseconds

// CreateController creates a new serial connection to a device
// and handles all communication and key interpretation
func CreateController(port *serial.SerialPort, kbArray [8]int, wg *sync.WaitGroup) {
	defer port.Close()
	// Creates Keyboard
	kb, err := keybd.NewKeyBonding()
	if err != nil {
		log.Fatal(err)
	}
	// Arduino and Keyboard Setup time
	time.Sleep(sleepTime * time.Millisecond)

	reader := bufio.NewReader(port)
	//Each HES key has it's on channel to signal completion
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
					wg.Done()
					return
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

// sendKeys handles sending keystrokes to host system
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
