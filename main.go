package main

import (
	"encoding/gob"
	"fmt"
	"github.com/nlopes/slack"
	"os"
	"sort"
    "strings"
)

// Get keys Array from Map
func Keys(assoc map[string]string) []string {
	keys := make([]string, len(assoc))
	i := 0

	for k := range assoc {
		keys[i] = k
		i++
	}

	return keys
}

// Save data to gob file
func Save(fileName string, data interface{}) error {
	db, err := os.OpenFile(fileName, os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	encoder := gob.NewEncoder(db)
	err = encoder.Encode(data)
	if err != nil {
		return err
	}

	return nil
}

// Load data from gob file
func Load(fileName string, data *sort.StringSlice) error {
	db, err := os.Open(fileName)
	if err != nil {
		return err
	}

	decoder := gob.NewDecoder(db)
	err = decoder.Decode(&data)
	if err != nil {
		return err
	}

	return nil
}

// String binary search. golang's SearchStrings may be broken
func BinSearch(xs []string, x string) int {
	l := 0
	h := len(xs) - 1
	for l <= h {
		m := l + (h-l)/2
		if xs[m] < x {
			l = m + 1
		} else {
			h = m - 1
		}
	}

	if l >= len(xs) || xs[l] != x {
		return -1
	}
	return l
}

// Check is file exist.
func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func PostMessageToSlack(slackApi *slack.Client, channel slack.Channel, text string) error {
	/// post message
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "1 Argument required Slack Api Token")
		return
	}
	channelName := "emoji"
	saveFileName := "emojis.gob"

	/// Change Values By Command Line Arguments
	if len(os.Args) > 2 {
		channelName = os.Args[2]
	}
	if len(os.Args) > 3 {
		saveFileName = os.Args[3]
	}

	slackApi := slack.New(os.Args[1])

	/// Get All Emojis
	emojis, err := slackApi.GetEmoji()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to Get Emoji List: ", err)
		return
	}

	/// Find Emoji Channel
	channels, err := slackApi.GetChannels(true) // don't search archived channel
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to Get Channels: ", err)
		return
	}
	var channel slack.Channel
	for _, c := range channels {
		if c.Name == channelName {
			channel = c
			break
		}
	}
	if channel.Name != channelName {
		fmt.Fprintln(os.Stderr, "Failed to Get Channel "+channelName)
		return
	}

	/// Get Keys and Sort
	var keys sort.StringSlice = Keys(emojis)
	keys = append(keys, "HELLO")
	keys.Sort()

	/// If save file is not exist, don't compare
	if !Exists(saveFileName) {
		fd, err := os.Create(saveFileName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		fd.Close()

		err = Save(saveFileName, keys)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		return
	}

	/// Restore Yesterday's Emojis
	var stored sort.StringSlice
	Load(saveFileName, &stored)

    newEmojis := make([]string)

	/// Find New (Non Exists in Yesterday) Emojis
	for _, name := range keys {
		if BinSearch(stored, name) == -1 {
            newEmojis = append(newEmojis, ":" + name + ":")
		}
	}

    /// Post New Emoji to emoji channel
    params := slack.NewPostMessageParameters()
    params.Username = "Today's New Emoji"
    params.IconEmoji = ":tada:"
    _, _, err := slackApi.PostMessage(channel.ID, strings.Join(newEmojis, " "), params) /// PostMessage function is changed now on github, to rewrite here when package updated
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
    }

	/// Store Today's Emojis
	Save(saveFileName, keys)
}
