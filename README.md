# Stack to Slack

This [Slack](https://slack.com/) bot monitors [StackExchange](https://stackexchange.com/) tags and automatically
publishes new questions in configured Slack channels.

It has been initially created to post all questions related to the [API Platform](https://api-platform.com) framework
in a dedicated channel of the Slack of [Les-Tilleuls.coop](https://les-tilleuls.coop) (the company behind the framework).

[![Go Report Card](https://goreportcard.com/badge/github.com/dunglas/stack2slack)](https://goreportcard.com/report/github.com/dunglas/stack2slack)

## Installing

This bot is written in [Go](https://golang.org/) (golang), you need a proper install of Go to compile it from sources.

1. [Create new Slack bot](https://my.slack.com/services/new/bot) and grab the generated API token
2. Clone this repository: `git clone https://github.com/dunglas/stack2slack.git`
3. Get the  dependencies `go get`
4. Compile the app: `go build`
5. Start the daemon: `DEBUG=1 SLACK_API_TOKEN=<your-API-token> TAG_TO_CHANNEL='{"stackoverflow-tag": "slack-channel"}' ./stack2slack`
6. Finally, you need to invite the bot in channels it will post: `/invite @bot-name`

## Credits

Written by [Kévin Dunglas](https://dunglas.fr).
Sponsored by [Les-Tilleuls.coop](https://les-tilleuls.coop).
