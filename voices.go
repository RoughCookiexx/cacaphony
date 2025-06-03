package main

import (
	"math/rand"
	"time"
)

// VoiceIDs holds available voice IDs
var VoiceIDs []string

func init() {
	// Seed random once
	VoiceIDs = []string{
		"Sje7TfONYrTXvC6J8LFX",
		"pL1ziwpsCAwVo2fcZC4U",
		"zxrpZKR8aSGU8OrkJzzu",
		"4tRn1lSkEn13EVTuqb0g",
		"RPdRfxxQOaNxn1LtRQqm",
		"ryn3WBvkCsp4dPZksMIf",
		"fThYUEUmlC2mx7fgTahR",
		"j7KV53NgP8U4LRS2k2Gs",
		"FmJ4FDkdrYIKzBTruTkV",
		"YVyp28LAMQfmx8iIH88U",
		"vAnM0Y3nRTffHvlgOQyj",	 
		"rJUGhpkqtiT41x4AdrRb",
		"nNnyUbsckxi5E8tk3BWv",
		"iHSDkze8smw4MQyIFKPZ",
		"S52z11GzOWpe02sIvCCW",
		"FVQMzxJGPUBtfz1Azdoy",
		"nzeAacJi50IvxcyDnMXa",
		"dDpKZ6xv1gpboV4okVbc",
		"Tn4bhLlhD26sndFn0Kgw",
		"IZA1V6HYiBphGehfV21Q",
	}

	rand.Seed(time.Now().UnixNano())
}

// SelectRandomVoiceID picks a random voice ID from VoiceIDs and removes it from the slice
func SelectRandomVoiceID() (string, bool) {
	n := len(VoiceIDs)
	if n == 0 {
		return "", false // no voice IDs left
	}

	// Pick random index
	idx := rand.Intn(n)
	voiceID := VoiceIDs[idx]

	// Remove the selected voiceID from the slice by swapping with last and slicing off
	VoiceIDs[idx] = VoiceIDs[n-1]
	VoiceIDs = VoiceIDs[:n-1]

	return voiceID, true

}
