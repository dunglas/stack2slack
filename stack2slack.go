package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/nlopes/slack"
)

const waitBetweenChecks = time.Minute * 5

func main() {
	slackApiToken := os.Getenv("SLACK_API_TOKEN")
	if "" == slackApiToken {
		log.Fatalln("The environment variable \"SLACK_API_TOKEN\" is not set.")
	}

	tagToChannel := os.Getenv("TAG_TO_CHANNEL")
	if "" == tagToChannel {
		log.Fatalln("The environment variable \"TAG_TO_CHANNEL\" is not set.")
	}

	var tagToChannelName map[string]string
	err := json.Unmarshal([]byte(tagToChannel), &tagToChannelName)
	if err != nil {
		log.Panicln("Unable to parse JSON data provided in  \"TAG_TO_CHANNEL\".")
		log.Fatal(err)
	}

	stackSite := os.Getenv("STACK_SITE")
	if stackSite == "" {
		stackSite = "stackoverflow"
	}

	runSlackClient(slackApiToken, stackSite, tagToChannelName, os.Getenv("DEBUG") == "1")
}

func runSlackClient(slackApiToken string, stackSite string, tagToChannelName map[string]string, debug bool) {
	api := slack.New(slackApiToken)
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)

	slack.SetLogger(logger)
	api.SetDebug(debug)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			if ev.ConnectionCount != 1 {
				log.Fatalln("This bot is already connected")
			}

			tagToChannelId := make(map[string]string, len(tagToChannelName))

		OUTER:
			for tagName, channelName := range tagToChannelName {
				for _, channel := range ev.Info.Channels {
					if channelName == channel.Name {
						tagToChannelId[tagName] = channel.ID

						continue OUTER
					}
				}

				log.Fatalf("The channel \"%s\" doesn't exist.", channelName)
			}

			watchStack(rtm, tagToChannelId, stackSite)
		}
	}
}

func watchStack(rtm *slack.RTM, tagToChannelId map[string]string, stackSite string) {
	lastCreationDate := 0

	tags := make([]string, len(tagToChannelId))
	i := 0
	for tag := range tagToChannelId {
		tags[i] = tag
		i++
	}

	baseUrl := fmt.Sprintf("https://api.stackexchange.com/2.2/search?order=desc&sort=creation&tagged=%s&site=%s", strings.Join(tags, ";"), stackSite)
	for {
		var url string
		if lastCreationDate > 0 {
			url = fmt.Sprintf("%s&min=%d", baseUrl, lastCreationDate)
		} else {
			url = baseUrl
		}

		resp, err := http.Get(url)
		if err != nil {
			log.Print(err)

			time.Sleep(waitBetweenChecks)
			continue
		}

		type Owner struct {
			DisplayName string `json:"display_name"`
		}

		type Item struct {
			Tags         []string `json:"tags"`
			Owner        Owner    `json:"owner"`
			CreationDate int      `json:"creation_date"`
			Title        string   `json:"title"`
			Link         string   `json:"link"`
		}

		type Response struct {
			Items []Item `json:"items"`
		}

		var stackResponse = new(Response)
		err = json.NewDecoder(resp.Body).Decode(&stackResponse)
		if err != nil {
			log.Print(err)

			time.Sleep(waitBetweenChecks)
			continue
		}

		for _, item := range stackResponse.Items {
			for _, tag := range item.Tags {
				if channelId, ok := tagToChannelId[tag]; ok {
					fmt.Printf("%v", channelId)

					rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("%s (%s) by %s. Tags: %s\n", item.Title, item.Link, item.Owner.DisplayName, strings.Join(item.Tags, ", ")), channelId))
				}
			}

			if item.CreationDate > lastCreationDate {
				lastCreationDate = item.CreationDate
			}
		}

		time.Sleep(waitBetweenChecks)
	}
}
