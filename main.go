package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"bytes"
	"strings"

	"github.com/gorilla/mux"
	"net/http"

	// Contains our user model and database logic
	"github.com/1egoman/slackpgp/users"
)

func main() {
	router := mux.NewRouter().StrictSlash(false)

	// Provide information when hitting the default route.
	router.HandleFunc("/", InfoHandler).Methods("GET")

	// Setting up a new user's public key
	router.HandleFunc("/onboard/{secret}", OnboardTemplateHandler).Methods("GET")
	router.HandleFunc("/onboard/{secret}", OnboardHandler).Methods("POST")
	router.HandleFunc("/onboard_success", UserEnteredPubKeyHandler).Methods("GET")

	// The main entrypoint - the slack webhook.
	router.HandleFunc("/webhook", WebhookHandler).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	fmt.Println("Listening on :" + port)
	http.ListenAndServe(":"+port, router)
}

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

func InfoHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, `<html>
    <head>
      <title>PGP Slackbot</title>
      <style>
        body {
          font-family: Helvetica, Arial, sans-serif;
          padding: 20px;
        }

        /* Hide sections that come in the future */
        #step-post-team-name, #step-post-slash-command, #step-complete {
          margin-top: 2em;
          margin-bottom: 1em;
          min-height: 200px;
          display: none;
        }

        button {
          display: block;
          font-size: 1.4em;
          background-color: #2ab27b;
          color: #ffffff;
          padding: 10px 20px;
          border: 0px;
          border-radius: 4px;
          margin-top: 1em;
          outline: none;
          cursor: pointer;
        }
        button:active {
          box-shadow: 0px 2px 5px rgba(0, 0, 0, 0.3) inset;
        }

        input[type=text] {
          padding: 2px 4px;
        }

        .warning {
          color: #EDB431;
        }

        code {
          padding: 2px 4px;
          font-size: 90%;
          color: #fff;
          background-color: #333;
          border-radius: 3px;
          -webkit-box-shadow: inset 0 -1px 0 rgba(0,0,0,.25);
          box-shadow: inset 0 -1px 0 rgba(0,0,0,.25);
        }

        li {
          padding: 2px 0px;
        }
      </style>
    </head>
    <body>
      <h1>Welcome to the PGP Slackbot!</h1>
      This slackbot provides a simple way to send encrypted messages with PGP.
      
      <span id="step-initial">
        <h2>Setup</h2>
        <p>
          First, what's your slack team url? <input type="text" id="slack-team-name" /><code>.slack.com</code>
          <button id="slack-team-continue">Continue</button>
        </p>
      </span>

      <span id="step-post-team-name">
        <h2>Setup slash command</h2>
        <p>
          <a id='slack-slash-command-setup' target="_blank" href="#dynamic-data-here">Click here</a> to setup a slash command for your team.
        </p>
        <p>
          Next to <em>Choose a Command</em>, type the command you'd like to use to talk to this bot. We like <code>/pgp</code>.
        </p>
        <span class="warning">Don't forget to click Add Slash Command Integration!</span>
        <button id="slack-slash-command-continue">Continue</button>
      </span>

      <span id="step-post-slash-command">
        <h2>Configure slash command</h2>
        Confirm the following textboxes are filled in:
        <ul>
          <li><em>URL</em> should read <code id="slack-command-hook-url">localhost:8000/webhook</code></li>
          <li><em>Method</em> should read <code>POST</code></li>
          <li><em>Escape channels, users, and links</em> should be <code>OFF</code></li>
          <li><small>(feel free to adjust the other settings to personalize the bot)</small></li>
        </ul>
        <span class="warning">Don't forget to click Save Integration!</span>
        <button id="slack-bot-created-continue">Continue</button>
      </span>

      <span id="step-complete">
        <h2>Test the bot</h2>
        <p>
          Now that the bot is set up, try running this command: <code>/pgp init</code>
        </p>
        <p>
          The slackbot will send you a link, click on the link to enter your public key.
        </p>
        <p>
          Finally, send an encrypted message to yourself: <code>/pgp @your-slack-username secret message!</code>
        </p>
      </span>

      <script>
      document.getElementById('slack-team-continue').onclick = function() {
        // When the user types in a name for their slack team, update the link to make a new
        // webhook.
        var teamName = document.getElementById('slack-team-name').value;
        document.getElementById('slack-slash-command-setup').href = "https://"+teamName+".slack.com/apps/new/A0F82E8CA"

        // Show the next step, and hide the previous step.
        document.getElementById('step-initial').style.display = 'none';
        document.getElementById('step-post-team-name').style.display = 'block';
      }

      document.getElementById('slack-slash-command-continue').onclick = function() {
        document.getElementById('slack-command-hook-url').innerHTML = location.origin + "/webhook";

        // Hide previous step, show next step.
        document.getElementById('step-post-team-name').style.display = 'none';
        document.getElementById('step-post-slash-command').style.display = 'block';
      }

      document.getElementById('slack-bot-created-continue').onclick = function() {
        // Hide previous step, show next step.
        document.getElementById('step-post-slash-command').style.display = 'none';
        document.getElementById('step-complete').style.display = 'block';
      }
      </script>
    </body>
  </html>`)
}

func UserEnteredPubKeyHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, `Cool, thanks for adding your public key!
Now, try to send someone else a message: /pgp @another-user my secret message!`)
}
