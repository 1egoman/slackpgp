package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"bytes"
	"strings"

  "net/http"

	// Contains our user model and database logic
	"github.com/1egoman/slackpgp/users"
)

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if len(r.Form["text"]) == 0 {
		io.WriteString(w, "No text arg!")
		return
	}
	if len(r.Form["user_name"]) == 0 {
		io.WriteString(w, "No user_name arg!")
		return
	}

	// If a slack token was specified, validate against that.
	if slack_token := os.Getenv("SLACK_VERIFICATION_TOKEN"); slack_token != "" && slack_token != r.Form["token"][0] {
		io.WriteString(w, "Invalid slack token. Please make sure the SLACK_VERIFICATION_TOKEN environment variable is set correctly.")
	}

	senderUsername := r.Form["user_name"][0]
	slashCommandText := r.Form["text"][0]
	slashCommandPayload := strings.Split(slashCommandText, " ")

	// Did the user type anything after the slash command?
	// ie, if the user just typed /pgp, then give them an error.
	if len(strings.Trim(strings.Join(slashCommandPayload, ""), " ")) == 0 {
		// Get the command that the user typed
		command := "/pgp"
		if len(r.Form["command"]) > 0 {
			command = r.Form["command"][0]
		}

		// Send them an error
		w.WriteHeader(404)
		io.WriteString(w, "Please specify a command (ie, `"+command+" init`) or send a message (ie, `"+command+" @user my secret message`)")
		return
	}

	switch r.Form["text"][0] {
	// Setup a new user.
	case "init":
		u := users.NewUser(senderUsername)
		u.EnableConfiguration()
		u.Create()

		// Where is the server hosted?
		hostname := os.Getenv("HOSTNAME")
		if hostname == "" {
			hostname = "http://localhost:8000"
		}

		// Give the user a path to set their public key.
		io.WriteString(w, "Click here to configure your public key: "+hostname+"/onboard/"+u.Secret)

	// Send an encrypted message to another user.
	default:
		// Remove '@' from start of username string
		recipientUsername := slashCommandPayload[0]
		if recipientUsername[0] == '@' {
			recipientUsername = recipientUsername[1:]
		}

		recipient, err := users.GetUserByUsername(recipientUsername)
		if recipient == nil {
			// Get the command that the user typed
			command := "/pgp"
			if len(r.Form["command"]) > 0 {
				command = strings.Split(r.Form["command"][0], " ")[0]
			}
			io.WriteString(w, "The user "+slashCommandPayload[0]+" doesn't exist or hasn't registered. Tell them to run `"+command+" init`.")
			return
		} else if err != nil {
			io.WriteString(w, err.Error())
		}

		// Send a placeholder response. This is so the slash command (with the private message!) won't
		// be shown to everyone.  Later, we'll send the actual message async (with the `response_url`
		// in the payload) which will be shown to everyone.
		io.WriteString(w, " ")

		// Encrypt a message.
		message := strings.Join(slashCommandPayload[1:], " ")
		encryptedMessageBody := map[string]interface{}{
			"response_type": "in_channel",
			"text": fmt.Sprintf(
				"Hey <@%s>, here's a message from <@%s>: \n ```%s```",
				recipientUsername,
				senderUsername,
				recipient.Encrypt(message),
			),
		}
		encodedEncryptedMessageBody, _ := json.Marshal(encryptedMessageBody)

		// Send encrypted message as an async message. See above on why this trick is required.
		_, err = http.Post(
			r.Form["response_url"][0],
			"application/json",
			bytes.NewBuffer([]byte(encodedEncryptedMessageBody)),
		)
	}
}
