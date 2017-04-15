package main

import (
  "fmt"
  "io"
  "os"
  "bytes"

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

	http.ListenAndServe(":8000", router)
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

  if err {
    w.WriteHeader(500)
    io.WriteString(w, err.Error())
  }
    // Respond.
    w.WriteHeader(201)
    io.WriteString(w, "Send message.")
  }
}
