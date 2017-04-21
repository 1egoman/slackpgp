package main

import (
	"fmt"
	"io"
	"os"

	"github.com/gorilla/mux"
	"net/http"
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




// When the user hits the root route, give them a onboarding guide.
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
