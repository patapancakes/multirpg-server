package main

import (
	"os"
	"strconv"
)

func getMapList() []uint16 {
	files, err := os.ReadDir(".")
	if err != nil {
		panic(err)
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

	return maps
}
