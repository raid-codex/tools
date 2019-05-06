package utils

import (
	"encoding/json"
	"log"
	"os"
)

func WriteToFile(filename string, val interface{}) error {
	log.Printf("writing to file %s\n", filename)
	f, errOpen := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if errOpen != nil {
		return errOpen
	}
	defer f.Close()
	data, errJSON := json.MarshalIndent(val, "", "  ")
	if errJSON != nil {
		return errJSON
	}
	nb, errWrite := f.Write(data)
	if errWrite != nil {
		return errWrite
	}
	log.Printf("wrote %d bytes to file\n", nb)
	return nil
}
