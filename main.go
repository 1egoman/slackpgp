package main

import (
  "fmt"
  "golang.org/x/net/context"

  "net/http"
  "github.com/gorilla/mux"

  slackbot "github.com/BeepBoopHQ/go-slackbot"
  "github.com/nlopes/slack"
)

func main() {
	// bot := slackbot.New(os.Getenv("SLACK_TOKEN"))
  //
  // // Send an encrypted message to somebody
	// bot.Hear("").MessageHandler(SendMessageHandler)
  //
  // // Give the bot your info so people can send encrypted messages to you.
	// bot.Messages(slackbot.DirectMessage).Subrouter().Hear("init").MessageHandler(OnboardingHandler)
	// bot.Run()

  router := mux.NewRouter().StrictSlash(false)

  // Setting up a new user's public key
  router.HandleFunc("/onboard/{secret}", OnboardTemplateHandler).Methods("GET")
  router.HandleFunc("/onboard/{secret}", OnboardHandler).Methods("POST")

	http.ListenAndServe(":8000", router)
}


func SendMessageHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
  fmt.Println("");
  fmt.Println("Text", evt.Text)
  fmt.Println("Message", evt)

  sender := User{SlackUser: "rgausnet", KeybaseUser: "rgausnet", PublicKey: publicKey}
  recipient := User{SlackUser: "rgausnet", KeybaseUser: "rgausnet", PublicKey: publicKey}

  // Encrypt a message to the recipient.
  encryptedMessage := recipient.Encrypt("Foo!")

  // Format the message to be sent via slack
  msg := fmt.Sprintf(
    "Hey <@%s>, here's a message from <@%s>: \n ```%s```",
    recipient.SlackUser,
    sender.SlackUser,
    encryptedMessage,
  )

  // Send the message
	bot.Reply(evt, msg, slackbot.WithTyping)
}

// When you DM the bot, it gives you a command
func OnboardingHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
  bot.Reply(evt, "Let's get set up. Click this link: http://localhost:8000/onboard/kyxmf34itxg4tay2", slackbot.WithTyping)
}


const publicKey = `-----BEGIN PGP PUBLIC KEY BLOCK-----
Version: Keybase OpenPGP v2.0.8
Comment: https://keybase.io/rgausnet

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
=+UMl
-----END PGP PUBLIC KEY BLOCK-----`
