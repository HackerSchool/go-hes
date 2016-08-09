package main

import (
    "io/ioutil"
    "strings"
    "github.com/Jguer/go-hes/driver"
    "log"
    "errors"
	"encoding/json"
    "sync"
    keybd "github.com/micmonay/keybd_event"
)

type keybinding struct {
	Selk  string `json:"select"`
	Start string `json:"start"`
	Up    string `json:"up"`
    Down  string `json:"down"`
    Left  string `json:"left"`
    Right string `json:"right"`
    A_key string `json:"a"`
    B_key string `json:"b"`
}

// findArduino looks for the file that represents the Arduino
// serial connection. Returns the fully qualified path to the
// device if we are able to find a likely candidate for an
// Arduino, otherwise an empty string if unable to find
// something that 'looks' like an Arduino device.
func findArduino() ([]string, int, error) {
    contents, _ := ioutil.ReadDir("/dev")
    n := 0
    var duinos []string
    // Look for what is mostly likely the Arduino device
    for _, f := range contents {
        if strings.Contains(f.Name(), "tty.usbserial") ||
        strings.Contains(f.Name(), "ttyUSB") {
            duinos = append(duinos, "/dev/" + f.Name())
            // log.Println("/dev/" + f.Name())
            n++
        }
    }

    if(n != 0) {
        return duinos, n, nil
    }

    // Have not been able to find a USB device that 'looks'
    // like an Arduino.
    return duinos, n, errors.New("Device Find: Unable to find HES")
}

// readProfile unmarshalls the json containing keybind information for HES
func readProfile(jsonStr []byte) ([]keybinding, error) {
    kbds := []keybinding{}
    var data map[string][]json.RawMessage
    err := json.Unmarshal(jsonStr, &data)
    if err != nil {
        log.Println(err)
        return nil, errors.New("Unable to Unmarshal json")
    }
    for _,profile := range data["keybindings"] {
        kbds = addKeybinding(profile, kbds)
    }
    return kbds, nil
}

// addKeybinding creates keybinding structs and appends them to a slice
func addKeybinding(profile json.RawMessage, kbds []keybinding) []keybinding {
    kb := keybinding{}
    if err := json.Unmarshal(profile, &kb); err != nil {
        log.Println(err)
    } else {
        if kb != *new(keybinding) {
            kbds = append(kbds, kb)
        }
    }
    return kbds
}

// translateKeybindings converts strings from keybinding struct to keybd identifiers
func translateKeybindings(kb keybinding) ([8]int) {
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

    var kb_array [8]int
    kb_array[0] = keymap[kb.Selk]
    kb_array[1] = keymap[kb.Start]
    kb_array[2] = keymap[kb.Up]
    kb_array[3] = keymap[kb.Down]
    kb_array[4] = keymap[kb.Left]
    kb_array[5] = keymap[kb.Right]
    kb_array[6] = keymap[kb.B_key]
    kb_array[7] = keymap[kb.A_key]

    return kb_array
}

func main() {
    var wg sync.WaitGroup
    // Find the device that represents the arduino serial
    // connection.
    duinos, n, err := findArduino();
    if err != nil {
        log.Fatal(err)
    }

    // Read the whole file at once
    // Should not be done but I'm feeling rather trusting of user input today
    b, err := ioutil.ReadFile("mappings.json")
    if err != nil {
        panic(err)
    }

    kbds, err := readProfile(b)
    if err != nil {
        log.Fatal(err)
    }

    for i:= 0; i < n; i++ {
        //log.Printf("%d %s %+v\n", i, duinos[i], kbds[i])
        wg.Add(1)
        go driver.CreateController(duinos[i], translateKeybindings(kbds[i]), &wg)
    }
    wg.Wait()
}

