package main

import (
	"encoding/json"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

type User struct {
	ID          []byte
	DisplayName string
	Name        string

	creds []webauthn.Credential
}

func (o *User) WebAuthnID() []byte {
	return o.ID
}

func (o *User) WebAuthnName() string {
	return o.Name
}

func (o *User) WebAuthnDisplayName() string {
	return o.DisplayName
}

func (o *User) WebAuthnIcon() string {
	return "https://pics.com/avatar.png"
}

func (o *User) WebAuthnCredentials() []webauthn.Credential {
	return o.creds
}

func (o *User) AddCredential(credential *webauthn.Credential) {
	o.creds = append(o.creds, *credential)
}

func (o *User) UpdateCredential(credential *webauthn.Credential) {
	for i, c := range o.creds {
		if string(c.ID) == string(credential.ID) {
			o.creds[i] = *credential
			return
		}
	}
}

// Convert User to DBUser
func (u *User) ToDBUser() *DBUser {
	return &DBUser{
		ID:          string(u.ID),
		DisplayName: u.DisplayName,
		Name:        u.Name,
	}
}

// Convert webauthn.Credential to DBCredential
func CredentialToDB(userID string, cred *webauthn.Credential) (*DBCredential, error) {
	// Serialize transport methods
	transportJSON, err := json.Marshal(cred.Transport)
	if err != nil {
		return nil, err
	}

	// Serialize authenticator data
	authData, err := json.Marshal(cred.Authenticator)
	if err != nil {
		return nil, err
	}

	// Serialize flags
	flagsData, err := json.Marshal(cred.Flags)
	if err != nil {
		return nil, err
	}

	return &DBCredential{
		ID:              string(cred.ID),
		UserID:          userID,
		PublicKey:       cred.PublicKey,
		AttestationType: cred.AttestationType,
		Transport:       string(transportJSON),
		Flags:           flagsData,
		Authenticator:   authData,
		SignCount:       cred.Authenticator.SignCount,
		CloneWarning:    cred.Authenticator.CloneWarning,
	}, nil
}

// Convert DBCredential to webauthn.Credential
func (dc *DBCredential) ToCredential() (*webauthn.Credential, error) {
	// Deserialize transport methods
	var transport []protocol.AuthenticatorTransport
	if err := json.Unmarshal([]byte(dc.Transport), &transport); err != nil {
		return nil, err
	}

	// Deserialize authenticator data
	var authenticator webauthn.Authenticator
	if err := json.Unmarshal(dc.Authenticator, &authenticator); err != nil {
		return nil, err
	}

	// Deserialize flags
	var flags webauthn.CredentialFlags
	if err := json.Unmarshal(dc.Flags, &flags); err != nil {
		return nil, err
	}

	return &webauthn.Credential{
		ID:              []byte(dc.ID),
		PublicKey:       dc.PublicKey,
		AttestationType: dc.AttestationType,
		Transport:       transport,
		Flags:           flags,
		Authenticator:   authenticator,
	}, nil
}

// Convert DBUser to User with credentials
func (du *DBUser) ToUser() (*User, error) {
	user := &User{
		ID:          []byte(du.ID),
		DisplayName: du.DisplayName,
		Name:        du.Name,
		creds:       make([]webauthn.Credential, 0),
	}

	// Convert credentials
	for _, dbCred := range du.Credentials {
		cred, err := dbCred.ToCredential()
		if err != nil {
			return nil, err
		}
		user.creds = append(user.creds, *cred)
	}

	return user, nil
}
