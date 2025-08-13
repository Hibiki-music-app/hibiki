package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
)

type PostgreSQLStore struct {
	log Logger
}

func NewPostgreSQLStore(log Logger) *PostgreSQLStore {
	return &PostgreSQLStore{
		log: log,
	}
}

func (p *PostgreSQLStore) GenSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (p *PostgreSQLStore) GetSession(token string) (webauthn.SessionData, bool) {
	p.log.Printf("[DEBUG] GetSession: %s", token)

	var dbSession DBSession
	db := GetDB()

	if err := db.Where("token = ? AND expires > ?", token, time.Now()).First(&dbSession).Error; err != nil {
		p.log.Printf("[DEBUG] Session not found or expired: %s", token)
		return webauthn.SessionData{}, false
	}

	// Deserialize session data
	var sessionData webauthn.SessionData
	if err := json.Unmarshal(dbSession.Data, &sessionData); err != nil {
		p.log.Printf("[ERROR] Failed to unmarshal session data: %v", err)
		return webauthn.SessionData{}, false
	}

	return sessionData, true
}

func (p *PostgreSQLStore) SaveSession(token string, data webauthn.SessionData) {
	p.log.Printf("[DEBUG] SaveSession: %s", token)

	db := GetDB()

	// Serialize session data
	sessionJSON, err := json.Marshal(data)
	if err != nil {
		p.log.Printf("[ERROR] Failed to marshal session data: %v", err)
		return
	}

	dbSession := DBSession{
		Token:   token,
		UserID:  string(data.UserID),
		Expires: data.Expires,
		Data:    sessionJSON,
	}

	// Upsert session
	if err := db.Save(&dbSession).Error; err != nil {
		p.log.Printf("[ERROR] Failed to save session: %v", err)
	}
}

func (p *PostgreSQLStore) DeleteSession(token string) {
	p.log.Printf("[DEBUG] DeleteSession: %s", token)

	db := GetDB()
	if err := db.Delete(&DBSession{}, "token = ?", token).Error; err != nil {
		p.log.Printf("[ERROR] Failed to delete session: %v", err)
	}
}

func (p *PostgreSQLStore) GetOrCreateUser(userName string) PasskeyUser {
	p.log.Printf("[DEBUG] GetOrCreateUser: %s", userName)

	db := GetDB()
	var dbUser DBUser

	// Try to get existing user with credentials
	err := db.Preload("Credentials").Where("name = ?", userName).First(&dbUser).Error
	if err != nil {
		p.log.Printf("[DEBUG] Creating new user: %s", userName)
		// Create new user
		dbUser = DBUser{
			ID:          userName,
			DisplayName: userName,
			Name:        userName,
		}

		if err := db.Create(&dbUser).Error; err != nil {
			p.log.Printf("[ERROR] Failed to create user: %v", err)
			// Return empty user as fallback
			return &User{
				ID:          []byte(userName),
				DisplayName: userName,
				Name:        userName,
			}
		}
	}

	// Convert to User
	user, err := dbUser.ToUser()
	if err != nil {
		p.log.Printf("[ERROR] Failed to convert DBUser to User: %v", err)
		// Return empty user as fallback
		return &User{
			ID:          []byte(userName),
			DisplayName: userName,
			Name:        userName,
		}
	}

	return user
}

func (p *PostgreSQLStore) SaveUser(user PasskeyUser) {
	p.log.Printf("[DEBUG] SaveUser: %s", user.WebAuthnName())

	db := GetDB()
	userName := user.WebAuthnName()

	// Convert to DBUser
	u := user.(*User) // Type assertion
	dbUser := u.ToDBUser()

	// Upsert user
	if err := db.Save(dbUser).Error; err != nil {
		p.log.Printf("[ERROR] Failed to save user: %v", err)
		return
	}

	// Delete existing credentials for this user
	if err := db.Delete(&DBCredential{}, "user_id = ?", userName).Error; err != nil {
		p.log.Printf("[ERROR] Failed to delete old credentials: %v", err)
		return
	}

	// Save new credentials
	for _, cred := range u.WebAuthnCredentials() {
		dbCred, err := CredentialToDB(userName, &cred)
		if err != nil {
			p.log.Printf("[ERROR] Failed to convert credential: %v", err)
			continue
		}

		if err := db.Create(dbCred).Error; err != nil {
			p.log.Printf("[ERROR] Failed to save credential: %v", err)
		}
	}
}

// CleanExpiredSessions removes expired sessions from database
func (p *PostgreSQLStore) CleanExpiredSessions() {
	db := GetDB()
	if err := db.Delete(&DBSession{}, "expires < ?", time.Now()).Error; err != nil {
		p.log.Printf("[ERROR] Failed to clean expired sessions: %v", err)
	}
}
