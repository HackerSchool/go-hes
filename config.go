package main

import (
	"encoding/json"
	keybd "github.com/Jguer/keybd_event"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

// readProfile unmarshalls the json containing keybind information for HES.
func readProfile() (kbds []keybinding, err error) {
	// Read the whole file at once
	// Should not be done but I'm feeling rather trusting of user input today
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	err = json.Unmarshal(raw, &kbds)
	if err != nil {
		return
	}

	return
}

func saveProfile(kbds []keybinding) (err error) {
	kbytes, err := json.Marshal(kbds)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(filename, kbytes, 0644)

	return
}

func configProfile(w http.ResponseWriter, r *http.Request) {
	kbds, err := readProfile()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v\n", kbds)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("config.html")
		t.Execute(w, kbds)
	} else {
		if err := r.ParseForm(); err != nil {
			log.Fatal(err)
		}
		// logic part of log in
		for key, values := range r.Form { // range over map
			for i, value := range values {
				kbds[i].populate(key, value)
			}
		}
		err = saveProfile(kbds)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (k *keybinding) populate(key string, value string) {
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
		"down":  keybd.VK_DOWN,
		"esc":   keybd.VK_ESC,
		"left":  keybd.VK_LEFT,
		"right": keybd.VK_RIGHT,
		"space": keybd.VK_ENTER,
		"up":    keybd.VK_UP,
	}

	if _, ok := keymap[value]; !ok {
		log.Println("Invalid Key")
		return
	}

	if key == "A" {
		k.A = value
	} else if key == "B" {
		k.B = value
	} else if key == "Start" {
		k.Start = value
	} else if key == "Select" {
		k.Select = value
	} else if key == "Left" {
		k.Left = value
	} else if key == "Right" {
		k.Right = value
	} else if key == "Up" {
		k.Up = value
	} else if key == "Down" {
		k.Down = value
	} else {
		log.Fatalln("Unable to populate keybinding")
	}
}

func startConfig() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", configProfile)
	log.Println("Config at http://127.0.0.1:9090/")
	err := http.ListenAndServe(":9090", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
