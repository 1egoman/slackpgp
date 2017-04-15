package main

import (
  "fmt"
  "io"
  "strings"

  "net/http"
  "github.com/gorilla/mux"

  "users"
)



// This is called on GET /onboard/{secret}. It renders the form for the user to set their public
// key.
func OnboardTemplateHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  user, err := users.GetUserBySecret(vars["secret"])
  fmt.Println("Found User", user)

  if err != nil {
    // On error, print out the error as a response
    w.WriteHeader(500)
    w.Header().Set("Content-Type", "text/plain")
    io.WriteString(w, err.Error())
  } else if user != nil && user.IsConfigurable {
    // We got a user that is configurable, so display a config page.
    w.Header().Set("Content-Type", "text/html")
    formattedTemplate := strings.Replace(ONBOARD_TEMPLATE, "{username}", user.Username, -1)
    io.WriteString(w, formattedTemplate)
  } else {
    // User isn't configurable. Respond with 404.
    w.WriteHeader(404)
    io.WriteString(w, "No such user.")
  }
}

// This posts the new public key to a user with the given secret. This also resets the
// `IsConfigurable` flag so the public key can't be reset uness the user initiates anotehr reset.
func OnboardHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  user, err := users.GetUserBySecret(vars["secret"])

  if err == nil && user != nil {
    // User with the specified secret key was found.

    // call to ParseForm makes form fields available.
    err := r.ParseForm()
    if err != nil {
      w.Header().Set("Content-Type", "text/plain")
      io.WriteString(w, err.Error())
    }

    // Update user public key, and block further configuration
    user.PublicKey = r.PostFormValue("key")
    user.IsConfigurable = false
    user.Save()

    // Redirect back to the root.
    http.Redirect(w, r, "/", 302)
  } else if err == nil && user == nil {
    // No user exists with that secret key
    w.WriteHeader(404)
    io.WriteString(w, "No such user.")
  } else {
    // Some other error occured.
    w.WriteHeader(500)
    w.Header().Set("Content-Type", "text/plain")
    io.WriteString(w, err.Error())
  }
}

const ONBOARD_TEMPLATE = `
<doctype html />
<html>
  <head>
    <title>Add your public key</title>
  </head>
  <body>
    <h1>Add your public key, {username}</h1>
    <form method="POST">
      <textarea
        name="key"
        placeholder="-----BEGIN PGP PUBLIC KEY BLOCK----- ..."
        style="width: 60em; height: 40em;"
      ></textarea>

      <br />
      <input type="submit" value="Set public key" />
    </form>
  </body>
</html>
`
