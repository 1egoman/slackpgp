package users

import (
  "os"
  "strings"

  "crypto/rand"
  "encoding/base64"

  "bytes"
  "golang.org/x/crypto/openpgp"
  "golang.org/x/crypto/openpgp/armor"

  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
  _ "github.com/jinzhu/gorm/dialects/postgres"
)

// How long should configuration secrets be in random characters before base64 encode?
const CONFIGURATION_SECRET_LENGTH = 32

var db *gorm.DB

type User struct {
  gorm.Model

  Username string
  PublicKey string

  // Specify if a user's public key can be configured. Also specify a secret that must be passed for
  // the configuration to happen.
  IsConfigurable bool
  Secret string
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
  db.Save(u)
}

func (u *User) Create() {
  db.Create(u)
}


// Given a public key and a message, encrypt the data with that pubkic key and return it.
// 
//   (input) => Encrypt => Armor => (output)
//
func (u *User) Encrypt(message string) string {
	encBuf := bytes.NewBuffer(nil)

  // Once the data has been encrypted, stream to the armorer
  armorer, err := armor.Encode(encBuf, "PGP MESSAGE", map[string]string{
    "Sent-By": "slackbot",
    "To-Slack-User": u.Username,
  })
	if err != nil {
		panic(err)
	}

  // Encrypt data from plaintext
  entityList, err := openpgp.ReadArmoredKeyRing(bytes.NewBufferString(u.PublicKey))
	encrypter, err := openpgp.Encrypt(armorer, entityList, nil, nil, nil)
	if err != nil {
		panic(err)
	}
	_, err = encrypter.Write([]byte(message))

	encrypter.Close()
	armorer.Close()

  return string(encBuf.Bytes())
}



func init() {
  var err error

  driver := os.Getenv("DATABASE_DRIVER")
  if driver == "" {
    panic("DATABASE_DRIVER is empty!")
  }
  path := os.Getenv("DATABASE_URL")
  if path == "" {
    panic("DATABASE_URL is empty!")
  }
  db, err = gorm.Open(driver, path)

  if err != nil {
    panic(err)
  }

  // Migrate user model
  db.AutoMigrate(&User{})
}


// Given a username as a string, return a pointer to a new user.
// THe user has public key configuration open.
func NewUser(username string) *User {
  user := User{Username: username}
  user.EnableConfiguration()
  return &user
}

// Given a secret, see if a configurable user can be found with that secret.
func GetUserBySecret(secret string) (*User, error) {
  user := &User{}
  err := db.Where("secret = ?", secret).First(user)
  if err.Error != nil {
    return nil, nil
  } else {
    return user, nil
  }
}

// Given a secret, see if a configurable user can be found with that secret.
func GetUserByUsername(username string) (*User, error) {
  user := &User{}
  err := db.Where("username = ?", username).First(user)
  if err.Error != nil {
    return nil, nil
  } else {
    return user, nil
  }
}
