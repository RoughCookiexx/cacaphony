package main

import (
	"math/rand"
	"time"

	"github.com/RoughCookiexx/gg_elevenlabs"
)

// VoiceIDs holds available voice IDs
var VoiceIDs []string

func init() {
	// Seed random once
	VoiceIDs, _ = gg_eleven.GetVoiceIDs()
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
