package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	
	"github.com/RoughCookiexx/gg_elevenlabs"
	"github.com/RoughCookiexx/gg_sse"
	"github.com/RoughCookiexx/gg_twitch_types"
	"github.com/RoughCookiexx/twitch_chat_subscriber"
)

var Users = make(map[string]string)

func handleMessage(message twitch_types.Message)(string) {
	const command = "!voice"
	if strings.Contains(message.Content, command) {	
		voiceURL := strings.TrimSpace(message.Content[len(command):])
		return setVoice(message.Tags.UserID, voiceURL)
	}
	
	return outburst(message.Content, message.Tags.UserID)

}

func outburst(message string, userName string)(string) {
	voiceId, exists := Users[userName]
	if !exists {
		voiceId, _ = SelectRandomVoiceID()
		Users[userName] = voiceId
		fmt.Printf("Assigned voice id %s to user %s", voiceId, userName)
	}	

	voiceResponse, err := gg_eleven.TextToSpeech(voiceId, message)
	if err != nil {
		fmt.Println(fmt.Sprintf("OH FUCK: %w", err))
		return "You fucked something up... or maybe I did?"
	}
	sse.SendBytes(voiceResponse)
	return ""
}

func setVoice(userId string, message string) (string) {
	elevenlabsUserId, voiceId, err := extractIDsFromURL(message)	
	if err != nil {
		fmt.Println("ERROR: ", err)
		return fmt.Sprintf("Could not add voice for some reason.. i dunno..")
	}

	existingVoiceId, exists := Users[userId]
	if !slices.Contains(VoiceIDs, voiceId) && existingVoiceId != voiceId {
		err = gg_eleven.AddSharedVoice(elevenlabsUserId, voiceId, userId)
		if err != nil {
			fmt.Println("ERROR: ", err)
			return fmt.Sprintf("Could not add voice for some reason.. i dunno..")
		}
	}

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
	port := 6972
	if len(os.Args) > 1 {
		for i, arg := range os.Args {
			if arg == "--port" && i+1 < len(os.Args) {
				p, err := strconv.Atoi(os.Args[i+1])
				if err == nil {
					port = p
				}
			}
		}
	}
	fmt.Println("Subscribing to chat messages")
	targetURL := "http://0.0.0.0:6969/subscribe"
	filterPattern := "PRIVMSG"
	twitch_chat_subscriber.SendRequestWithCallbackAndRegex(targetURL, handleMessage, filterPattern, port)
        http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
                w.WriteHeader(http.StatusOK)
                fmt.Fprintf(w, "OK")
        })
	sse.Start()
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
