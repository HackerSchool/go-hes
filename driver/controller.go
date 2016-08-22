package driver

import (
	"bufio"
	keybd "github.com/micmonay/keybd_event"
	"github.com/tarm/serial"
	"log"
	"strconv"
	"sync"
	"time"
)

const challenge string = "Hi. Who are you?"
const repeatDelay time.Duration = 160 //Milliseconds
const sleepTime time.Duration = 1500  //Milliseconds

// CreateController creates a new serial connection to a device
// and handles all communication and key interpretation
func CreateController(device string, kbArray [8]int, wg *sync.WaitGroup) {
	c := &serial.Config{Name: device, Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()
	// Creates Keyboard
	kb, err := keybd.NewKeyBonding()
	if err != nil {
		log.Fatal(err)
	}
	// Arduino and Keyboard Setup time
	time.Sleep(sleepTime * time.Millisecond)

	n, err := s.Write([]byte(challenge))
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 128)
	n, err = s.Read(buf)
	if err != nil {
		log.Fatal(err)
	} else if n == 0 {
		log.Fatal("No response from HES")
	}

	log.Println("Communication established with " + device)
	reader := bufio.NewReader(s)

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

// getKey converts the received key to a keyboard key code
func getKey(key byte, kbArray [8]int) int {
	x, _ := strconv.Atoi(string(key))
	return kbArray[x]
}

// sendKeys handles sending keystrokes to host system
func sendKeys(signal chan bool, key byte, kb keybd.KeyBonding, kbArray [8]int) {
	kb.SetKeys(getKey(key, kbArray)) //set keys
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
