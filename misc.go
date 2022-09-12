package main

import "math/rand"

func generateLobbyCode() []byte {
	const runes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const runesLen = len(runes)

	b := make([]byte, 6)
    for i := range b {
        b[i] = runes[rand.Intn(runesLen)]
    }

    return b
}