package utils

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
)

func WriteToFile(filename string, val interface{}) error {
	f, errOpen := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if errOpen != nil {
		return errOpen
	}
	defer f.Close()
	var toWrite []byte
	if _, ok := val.([]byte); ok {
		toWrite = val.([]byte)
	} else {
		buf := bytes.NewBufferString("")
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "  ")
		errJSON := enc.Encode(val)
		if errJSON != nil {
			return errJSON
		}
		toWrite = buf.Bytes()
	}
	nbr, errWrite := f.Write(toWrite)
	if errWrite != nil {
		return errWrite
	}
	log.Printf("writing %d bytes to file %s\n", nbr, filename)
	return nil
}
