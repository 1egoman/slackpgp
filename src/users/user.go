package users

import (
  // "fmt"
  // "os"
  // "database/sql"
  // _ "github.com/mattn/go-sqlite3"
  "crypto/rand"
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

var users []User

func init() {
  users = append(users, User{
    Username: "rgausnet",
    PublicKey: "pubkey",
    Secret: "abcdef",
    IsConfigurable: true,
  })
}

// var db *sql.DB
// func initDatabase() {
//   var err error
//
//   // Where should the database be located?
//   DATABASE_NAME := "./foo.db"
//   if a := os.Getenv("DATABASE_NAME"); a != "" {
//     DATABASE_NAME = a
//   }
//
//   // Open the database. `db` is a global used by all user methods.
//   db, _ = sql.Open("sqlite3", DATABASE_NAME)
//   // defer db.Close()
//
//   // Create the users table if it doesn't exist.
// 	sqlStmt := `
// 	CREATE TABLE IF NOT EXISTS users (
//     id integer not null primary key,
//     username text
//     pubkey text
//     isconfig bool
//     secret text
//   );
// 	DELETE FROM users;
// 	`
//   _, err = db.Exec(sqlStmt)
// 	if err != nil {
//     panic(err)
// 	}
//
//   fmt.Printf("* Created database in %s\n", DATABASE_NAME)
// }

// Given a sql row, unpack it into a struct
// func usersRowToStruct(row *sql.Row) (*User, error) {
//   var user User
//   err := row.Scan(&user.Username, &user.PublicKey, &user.IsConfigurable, &user.ConfigurationSecret)
//   if err != nil {
//     return nil, err
//   } else {
//     return &user, nil
//   }
// }

// Set a user to be configurable.
func (u *User) EnableConfiguration() error {
  configSecret := make([]byte, CONFIGURATION_SECRET_LENGTH)
  _, err := rand.Read(configSecret)

  if err != nil {
    return err;
  } else {
    u.IsConfigurable = true
    u.Secret = string(configSecret);
    return nil;
  }
}



func NewUser(username string) *User {
  return &User{Username: username}
}

func NewUserWithKey(username string, publicKey string) *User {
  return &User{Username: username, PublicKey: publicKey}
}



func GetUserBySecret(secret string) (*User, error) {
  for _, user := range users {
    if user.Secret == secret {
      return &user, nil
    }
  }
  return nil, nil
}
