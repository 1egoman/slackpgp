package main

import (
  // "fmt"
  "bytes"
  "golang.org/x/crypto/openpgp"
  "golang.org/x/crypto/openpgp/armor"
)

type User struct {
  SlackUser string
  KeybaseUser string
  PublicKey string
}


// Given a public key and a message, encrypt the data with that pubkic key and return it.
// 
//   (input) => Encrypt => Armor => (output)
//
func (u User) Encrypt(message string) string {
	encbuf := bytes.NewBuffer(nil)

  // Once the data has been encrypted, stream to the armorer
  w, err := armor.Encode(encbuf, "PGP MESSAGE", map[string]string{
    "Sent-By": "slackbot",
    "To-Keybase-User": u.KeybaseUser,
    "To-Slack-User": u.SlackUser,
  })
	if err != nil {
		panic(err)
	}

  // Encrypt data from plaintext
  entityList, err := openpgp.ReadArmoredKeyRing(bytes.NewBufferString(u.PublicKey))
	plaintext, err := openpgp.Encrypt(w, entityList, nil, nil, nil)
	if err != nil {
		panic(err)
	}
	_, err = plaintext.Write([]byte(message))

	plaintext.Close()
	w.Close()

  return string(encbuf.Bytes())
}
