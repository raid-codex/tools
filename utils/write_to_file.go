package utils

import (
	"encoding/json"
	"os"
)

func WriteToFile(filename string, val interface{}) error {
	f, errOpen := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if errOpen != nil {
		return errOpen
	}
	defer f.Close()
	data, errJSON := json.MarshalIndent(val, "", "  ")
	if errJSON != nil {
		return errJSON
	}
	_, errWrite := f.Write(data)
	if errWrite != nil {
		return errWrite
	}
	return nil
}
