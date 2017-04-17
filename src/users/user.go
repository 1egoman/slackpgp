package users

import (
  "fmt"
  "strings"

  "crypto/rand"
  "encoding/base64"

  "bytes"
  "golang.org/x/crypto/openpgp"
  "golang.org/x/crypto/openpgp/armor"
)

// How long should configuration secrets be in random characters before base64 encode?
const CONFIGURATION_SECRET_LENGTH = 32

type User struct {
  Username string
  PublicKey string

  // Specify if a user's public key can be configured. Also specify a secret that must be passed for
  // the configuration to happen.
  IsConfigurable bool
  Secret string
}

var users []*User

// Set a user to be configurable by setting the boolean and generating a new secret.
func (u *User) EnableConfiguration() error {
  configSecret := make([]byte, CONFIGURATION_SECRET_LENGTH)
  _, err := rand.Read(configSecret)

  if err != nil {
    return err;
  } else {
    u.IsConfigurable = true

    // Encode secret to base64. Replace plus and slash with url safe characters.
    u.Secret = base64.StdEncoding.EncodeToString(configSecret)
    u.Secret = strings.Replace(u.Secret, "/", "_", -1)
    u.Secret = strings.Replace(u.Secret, "+", "-", -1)
    return nil;
  }
}

// When updates are made to a user struct, this will sync any changes back to disk.
func (u *User) Save() {
  for ct, user := range users {
    if u.Username == user.Username {
      fmt.Println("Saving...", u)
      users[ct] = u
      return
    }
  }

  // If the user isn't in the list, then add the user.
  users = append(users, u)
}


// Given a public key and a message, encrypt the data with that pubkic key and return it.
// 
//   (input) => Encrypt => Armor => (output)
//
func (u *User) Encrypt(message string) string {
	encbuf := bytes.NewBuffer(nil)

  // Once the data has been encrypted, stream to the armorer
  w, err := armor.Encode(encbuf, "PGP MESSAGE", map[string]string{
    "Sent-By": "slackbot",
    "To-Slack-User": u.Username,
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




func NewUser(username string) *User {
  user := NewUser(username)
  user.EnableConfiguration()
  return user
}

// Given a secret, see if a configurable user can be found with that secret.
func GetUserBySecret(secret string) (*User, error) {
  for _, user := range users {
    if user.Secret == secret && user.IsConfigurable {
      return user, nil
    }
  }

  // No user found! :/
  return nil, nil
}

// Given a secret, see if a configurable user can be found with that secret.
func GetUserByUsername(username string) (*User, error) {
  for _, user := range users {
    if user.Username == username {
      return user, nil
    }
  }

  // No user found! :/
  return nil, nil
}
