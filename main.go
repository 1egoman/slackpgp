package main

import (
  "fmt"
  "io"
  "encoding/json"

  "bytes"
  "strings"

  "net/http"
  "github.com/gorilla/mux"

  // Contains our user model and database logic
  "users"
)

func main() {
  router := mux.NewRouter().StrictSlash(false)

  // Setting up a new user's public key
  router.HandleFunc("/onboard/{secret}", OnboardTemplateHandler).Methods("GET")
  router.HandleFunc("/onboard/{secret}", OnboardHandler).Methods("POST")

  // The main entrypoint - the slack webhook.
  router.HandleFunc("/webhook", WebhookHandler).Methods("POST")

	http.ListenAndServe(":8000", router)
}

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()

  senderUsername := r.Form["user_name"][0]
  slashCommandText := r.Form["text"][0]
  slashCommandPayload := strings.Split(slashCommandText, " ")

  switch r.Form["text"][0] {

    // Setup a new user.
    case "init":
      u := users.NewUser(senderUsername)
      u.EnableConfiguration()
      u.Save()

      io.WriteString(w, "Click here to set your pgp key: http://localhost:8000/onboard/"+u.Secret)


    // Send an encrypted message to another user.
    default:
      // Remove '@' from start of username string
      recipientUsername := slashCommandPayload[0]
      if recipientUsername[0] == '@' {
        recipientUsername = recipientUsername[1:]
      }

      recipient, err := users.GetUserByUsername(recipientUsername)
      if recipient == nil {
        io.WriteString(w, "The user "+slashCommandPayload[0]+" doesn't exist or hasn't registered. Tell them to run `/pgp init`.")
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
