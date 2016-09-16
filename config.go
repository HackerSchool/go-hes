package main

import (
	"encoding/json"
	"io/ioutil"
)

// readProfile unmarshalls the json containing keybind information for HES.
func readProfile() (kbds []keybinding, err error) {
	// Read the whole file at once
	// Should not be done but I'm feeling rather trusting of user input today
	raw, err := ioutil.ReadFile("mappings.json")
	if err != nil {
		return
	}

	err = json.Unmarshal(raw, &kbds)
	if err != nil {
		return
	}

	return
}
