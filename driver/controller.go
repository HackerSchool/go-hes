package driver

import (
	"bufio"
	"fmt"
	"runtime"
	"strconv"
	"time"

	keybd "github.com/jguer/keybd_event"
	s "go.bug.st/serial.v1"
)

const sleepTime time.Duration = 1 //Milliseconds

// CreateController reads all valid serial ports
// and handles all communication and key interpretation.
func CreateController(port s.Port, kbArray [8]int, exit chan<- bool, disconnect chan<- bool) {
	keybd.NewKeyBonding()
	// Arduino and Keyboard Setup time
	if runtime.GOOS == "linux" {
		time.Sleep(sleepTime * time.Second)
	}

	reader := bufio.NewReader(port)
	//Each HES key has its on channel to signal completion
	index := 0
	counter := 0

	status := make(chan int)

	go sendKeys(status, kbArray)

	for {
		readKey, err := reader.ReadBytes('\n')
		// fmt.Println("Read keys", string(readKey))
		if err != nil {
			disconnect <- true
			fmt.Println("Lost connection to one controller. Driver Reset.")
			return
		}

		if readKey[0] == 'P' {
			index, _ = strconv.Atoi(string(readKey[1]))
			status <- index

			//Pressing Select 5 Times closes the controller
			if index == 0 {
				counter++
				if counter == 5 {
					exit <- true
				}
			} else {
				counter = 0
			}

		} else if readKey[0] == 'R' {
			index, _ = strconv.Atoi(string(readKey[1]))
			status <- index
		}
	}
}

// sendKeys handles sending keystrokes to host system.
func sendKeys(signal <-chan int, kbArray [8]int) {
	pressed := [8]bool{}
	var i int
	for {
		i = <-signal
		keybdn := kbArray[i] //set keys
		if !pressed[i] {
			pressed[i] = true
			keybd.DownKey(kbArray[i])

		} else {
			pressed[i] = false
			keybd.UpKey(keybdn)
		}

		keybd.Sync()
	}
}
