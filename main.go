package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	
	"github.com/RoughCookiexx/gg_elevenlabs"
	"github.com/RoughCookiexx/gg_sse"
	"github.com/RoughCookiexx/twitch_chat_subscriber"
)

var Users = make(map[string]string)

func handleMessage(message string)(string) {
	userName, err := getUserName(message)
	if err != nil {
		fmt.Println("Display name not found.")
		return ""
	}

	if strings.Contains(message, "!voice") {	
		message = afterLastChar(message, " ")
		return setVoice(userName, message)
	}
	message = afterLastChar(message, ":")
	
	return outburst(message, userName)

}

func outburst(message string, userName string)(string) {
	voiceId, exists := Users[userName]
	if !exists {
		voiceId, _ = SelectRandomVoiceID()
		Users[userName] = voiceId
		fmt.Printf("Assigned voice id %s to user %s", voiceId, userName)
	}	

	voiceResponse := gg_eleven.TextToSpeech(voiceId, message)
	sse.SendBytes(voiceResponse)
	return ""
}

func setVoice(userId string, message string) (string) {
	elevenlabsUserId, voiceId, err := extractIDsFromURL(message)	
	if err != nil {
		fmt.Println("ERROR: ", err)
		return fmt.Sprintf("Could not add voice for some reason.. i dunno..")
	}

	err = gg_eleven.AddSharedVoice(elevenlabsUserId, voiceId, userId)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return fmt.Sprintf("Could not add voice for some reason.. i dunno..")
	}

	existingVoiceId, exists := Users[userId]
	Users[userId] = voiceId
	if !exists {
		VoiceIDs = append(VoiceIDs, existingVoiceId)
	}	
	return fmt.Sprintf("Set voice %s for %s", voiceId, userId)
}

func getUserName(msg string) (string, error) {
	re := regexp.MustCompile(`display-name=([^;]+)`)
	match := re.FindStringSubmatch(msg)
	if len(match) > 1 {
		return match[1], nil
	} else {
		return "", errors.New("Display name not found")	
	}
}

func afterLastChar(s string, c string) string {
	idx := strings.LastIndex(s, c)
	if idx == -1 || idx+1 >= len(s) {
		return ""
	}
	return s[idx+1:]
}     

func extractIDsFromURL(url string) (string, string, error) {
	success, url := trimAndValidateURL(url)
	if !success {
		return "", "", fmt.Errorf("the URL is malformed, you idiot")
	}
	parts := strings.Split(strings.Trim(url, "/"), "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("the URL is malformed, you idiot")
	}
	
	n := len(parts)
	userID := parts[n-2]
	voiceID := parts[n-1]
	
	return userID, voiceID, nil
}

func trimAndValidateURL(message string) (bool, string) {
	trimmedMsg := strings.TrimSpace(message)
	
	u, err := url.Parse(trimmedMsg)
	if err != nil || u.Scheme == "" || u.Host == "" {
		fmt.Printf("Failed to get URL from: %s", message)
		return false, trimmedMsg
	}
	
	return true, trimmedMsg
}

func main() {
	fmt.Println("Subscribing to chat messages")
	targetURL := "http://0.0.0.0:6969/subscribe"
	filterPattern := "PRIVMSG"
	twitch_chat_subscriber.SendRequestWithCallbackAndRegex(targetURL, handleMessage, filterPattern, 6972)
	sse.Start()
	http.ListenAndServe((":6972"), nil)
}
