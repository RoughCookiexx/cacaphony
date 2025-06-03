package main

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	
	"github.com/RoughCookiexx/gg_elevenlabs"
	"github.com/RoughCookiexx/gg_sse"
	"github.com/RoughCookiexx/twitch_chat_subscriber"
)

var Users = make(map[string]string)

func outburst(message string)(string) {
	userName, err := getUserName(message)
	if err != nil {
		fmt.Println("Display name not found.")
		return ""
	}
	message = afterLastColon(message)
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

func getUserName(msg string) (string, error) {
	re := regexp.MustCompile(`display-name=([^;]+)`)
	match := re.FindStringSubmatch(msg)
	if len(match) > 1 {
		return match[1], nil
	} else {
		return "", errors.New("Display name not found")	
	}
}

func afterLastColon(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx == -1 || idx+1 >= len(s) {
		return ""
	}
	return s[idx+1:]
}     

func main() {
	fmt.Println("Subscribing to chat messages")
	targetURL := "http://0.0.0.0:6969/subscribe"
	filterPattern := "PRIVMSG"
	twitch_chat_subscriber.SendRequestWithCallbackAndRegex(targetURL, outburst, filterPattern, 6972)
	sse.Start()
	http.ListenAndServe((":6972"), nil)
}
