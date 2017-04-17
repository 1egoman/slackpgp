package main

import (
  "fmt"
  "io"
  "os"

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

  router.HandleFunc("/encrypt/{username}", EncryptionHandler).Methods("POST")

  router.HandleFunc("/webhook", WebhookHandler).Methods("POST")

	http.ListenAndServe(":8000", router)
}

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()

  // body, _ := ioutil.ReadAll(r.Body);
  fmt.Println(r.Form)
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

      // Send a placeholder response. This is so the slash command won't be shown to everyone.
      // Later, we'll send the actual message async (with the `response_url` in the payload) which
      // will be shown to everyone.
      io.WriteString(w, "...")

      // Encrypt a message.
      message := strings.Join(slashCommandPayload[1:], " ")
      encryptedMessageBody := fmt.Sprintf(`{
        "response_type": "in_channel",
        "text": "Hey <@%s>, here's a message from <@%s>:\n",
        "attachments": [
            {
              "text": "%s"
            }
        ]
      }`, recipientUsername, senderUsername, recipient.Encrypt(message))

      // Send as an async message. See above on why this trick is required.
      _, err = http.Post(
        r.Form["response_url"][0],
        "application/json",
        bytes.NewBuffer([]byte(encryptedMessageBody)),
      )
  }
}

func EncryptionHandler(w http.ResponseWriter, r *http.Request) {
  // Get the recipient of the message from the url, and turn the recipient username
  // into a struct.
  vars := mux.Vars(r)
  recipient, err := users.GetUserByUsername(vars["username"])
  if err != nil {
    panic(err)
  } else if recipient == nil {
    w.WriteHeader(404)
    w.Header().Set("Content-Type", "text/plain")
    io.WriteString(w, "No such user with that username exists.")
  } else if len(recipient.PublicKey) == 0 {
    w.WriteHeader(400)
    w.Header().Set("Content-Type", "text/plain")
    io.WriteString(w, "Recipient user doesn't have a public key defined.")
  }

  // Encrypt a message to the recipient.
  encryptedMessage := recipient.Encrypt("Foo!")

  // Format the message to be sent via slack
  msg := fmt.Sprintf(
    "{\"text\": \"Hey <@%s>, here's a message: \n ```%s```\"}",
    recipient.Username,
    encryptedMessage,
  )

  // Send the message
  _, err = http.Post(
    os.Getenv("SLACK_INCOMING_WEBHOOK_URL"),
    "application/json",
    bytes.NewBuffer([]byte(msg)),
  )

  if err != nil {
    w.WriteHeader(500)
    io.WriteString(w, err.Error())
  } else {
    // Respond.
    w.WriteHeader(201)
    io.WriteString(w, "Sent message.")
  }
}
