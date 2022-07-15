package main

/*
multirpg-server
https://github.com/Gamizard/multirpg-server

Copyright (C) 2022 azarashi <azarashi@majestaria.fun>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.
*/

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
