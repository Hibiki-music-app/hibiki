package repository

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/Hibiki-music-app/hibiki/backend/internal"
	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"
)

type UserWebAuthn struct {
	ID          []byte
	DisplayName string
	Name        string

	creds []webauthn.Credential
}

func (o *UserWebAuthn) WebAuthnID() []byte {
	return o.ID
}

func (o *UserWebAuthn) WebAuthnName() string {
	return o.Name
}

func (o *UserWebAuthn) WebAuthnDisplayName() string {
	return o.DisplayName
}

func (o *UserWebAuthn) WebAuthnCredentials() []webauthn.Credential {
	return o.creds
}

func (o *UserWebAuthn) AddCredential(credential *webauthn.Credential) {
	o.creds = append(o.creds, *credential)
}

func (o *UserWebAuthn) UpdateCredential(credential *webauthn.Credential) {
	for i, c := range o.creds {
		if string(c.ID) == string(credential.ID) {
			o.creds[i] = *credential
		}
	}
}

func GetOrCreateUser(username string) (*User, error) {
	db := internal.GetDB()
	var user User
	result := db.First(&user, "name = ?", username)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		user = User{
			Name:        username,
			DisplayName: username,
		}
		if err := db.Create(&user).Error; err != nil {
			return nil, err
		}
	} else if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func GenSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func SaveSession(token string, data webauthn.SessionData) {
	db := internal.GetDB()
	var session Session
	session.SessionKey = data
}

func GetSession()
