package gohome

import (
	"encoding/base64"
	"math/rand"

	"github.com/markdaws/gohome/validation"
	errExt "github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user of the system
type User struct {
	ID        string
	Login     string
	Prefs     UserPrefs
	HashedPwd string
	Salt      string

	// TODO: right now users are hidden but here so that the
	// app architecture supports them, once they are exposed
	// we will need to add hashed password values, salts etc
}

// Validate verifies the user object is in a good state
func (u *User) Validate() *validation.Errors {
	errors := &validation.Errors{}

	if u.Login == "" {
		errors.Add("required field", "Login")
	}

	if errors.Has() {
		return errors
	}
	return nil
}

// SetPassword sets the users password
func (u *User) SetPassword(pwd string) error {
	errors := &validation.Errors{}

	if pwd == "" {
		errors.Add("required field", "Password")
	}

	if errors.Has() {
		return errors
	}

	salt, err := u.generateSalt()
	if err != nil {
		return errExt.Wrap(err, "failed to generate user salt")
	}

	u.Salt = salt
	hashed, err := bcrypt.GenerateFromPassword([]byte(pwd+u.Salt), bcrypt.DefaultCost)
	if err != nil {
		return errExt.Wrap(err, "failed to hash password")
	}
	u.HashedPwd = string(hashed)
	return nil
}

// VerifyPassword returns nil if the password matches the users password, otherwise an error is returned
func (u *User) VerifyPassword(pwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.HashedPwd), []byte(pwd+u.Salt))
}

// generateSalt generates a unique salt we can use in the password hashing
func (u *User) generateSalt() (string, error) {

	b := make([]byte, 128)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

// UserPrefs contains all of the user specific preferences
type UserPrefs struct {
	// UI are user UI preference settings
	UI UIPrefs
}

// UIPrefs contains preferences for the UI
type UIPrefs struct {
	// HiddenZones is a map keyed by zone IDs for zones that should
	// not be displayed in the UI.
	HiddenZones map[string]bool
}
