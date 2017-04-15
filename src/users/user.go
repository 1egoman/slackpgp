package users

import (
  "fmt"
  "strings"
  // "os"
  // "database/sql"
  // _ "github.com/mattn/go-sqlite3"
  "crypto/rand"
  "encoding/base64"
)

// How long should configuration secrets be in length?
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

func init() {
  user := NewUser("rgausnet")
  fmt.Println("Created User: ", user)

  // Create a test user
  users = append(users, user)
}

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
}



func NewUser(username string) *User {
  user := User{Username: username, IsConfigurable: true}
  user.EnableConfiguration()
  return &user
}


// Given a secret, see if a configurable user can be found with that secret.
func GetUserBySecret(secret string) (*User, error) {
  for _, user := range users {
    if user.Secret == secret && user.IsConfigurable {
      return user, nil
    }
  }
  return nil, nil
}
