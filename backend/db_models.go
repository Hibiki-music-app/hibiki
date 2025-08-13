package main

import (
	"time"
)

// DBUser represents a user in the database
type DBUser struct {
	ID          string `gorm:"primaryKey"`
	DisplayName string
	Name        string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Credentials []DBCredential `gorm:"foreignKey:UserID"`
}

// DBCredential represents a WebAuthn credential in the database
type DBCredential struct {
	ID              string `gorm:"primaryKey"`
	UserID          string `gorm:"index"`
	PublicKey       []byte
	AttestationType string
	Transport       string // JSON array of transport methods
	Flags           []byte // User presence, verification flags etc.
	Authenticator   []byte // JSON serialized authenticator data
	CreatedAt       time.Time
	UpdatedAt       time.Time
	SignCount       uint32
	CloneWarning    bool
}

// DBSession represents a session in the database
type DBSession struct {
	Token     string `gorm:"primaryKey"`
	UserID    string
	Challenge string
	Expires   time.Time
	CreatedAt time.Time
	Data      []byte // JSON serialized SessionData
}

// TableName specifies the table name for DBUser
func (DBUser) TableName() string {
	return "users"
}

// TableName specifies the table name for DBCredential
func (DBCredential) TableName() string {
	return "credentials"
}

// TableName specifies the table name for DBSession
func (DBSession) TableName() string {
	return "sessions"
}
