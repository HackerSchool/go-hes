package main

import (
	"encoding/json"
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

	if r.Method == "GET" {
		t, _ := template.ParseFiles("tmpl/login.html")
		t.Execute(w, kbds)
	} else {
		r.ParseForm()
		// logic part of log in
		for key, values := range r.Form { // range over map
			log.Println(key, values[0])
		}
	}
}

func startConfig() {
	http.HandleFunc("/login", configProfile)
	log.Println("Test server at http://127.0.0.1:9090/login")
	err := http.ListenAndServe(":9090", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
