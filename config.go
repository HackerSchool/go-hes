package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
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

		log.Printf("%+v\n", kbds)
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
		log.Println("Saved Config. Shutting down program.")
		time.Sleep(time.Second * 10)
		os.Exit(0)
	}
}

func (k *keybinding) populate(key string, value string) {
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
