package main

import (
	"os"
	"path/filepath"
	"strconv"
)

func getMapList() ([]uint16, error) {
	files, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}

	var maps []uint16
	for _, file := range files {
		if len(file.Name()) == 11 && file.Name()[7:] == ".lmu" {
			id, err := strconv.Atoi(file.Name()[3:7])
			if err != nil {
				panic(err)
			}

			maps = append(maps, uint16(id))
		}
	}

	return maps, nil
}

func getCharSetList() []string {
	files, err := os.ReadDir("CharSet")
	if err != nil {
		panic(err)
	}

	var charsets []string
	for _, file := range files {
		charsets = append(charsets, file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))])
	}

	return charsets
}

/*func getSoundList() []string {
	files, err := os.ReadDir("Sound")
	if err != nil {
		panic(err)
	}

	var sounds []string
	for _, file := range files {
		sounds = append(sounds, file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))])
	}

	return sounds
}

func getSystemList() []string {
	files, err := os.ReadDir("System")
	if err != nil {
		panic(err)
	}

	var systems []string
	for _, file := range files {
		systems = append(systems, file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))])
	}

	return systems
}*/

func isValidCharSet(charset string) bool {
	for _, char := range getCharSetList() {
		if char == charset {
			return true
		}
	}

	return false
}
