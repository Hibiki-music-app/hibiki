package repository

import "time"

type User struct {
	ID          uint      `gorm:"primaryKey"`
	Username    string    `gorm:"type:text;not null;unique"`
	DisplayName string    `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

type Credential struct {
	ID              uint      `gorm:"primaryKey"`
	UserID          uint      `gorm:"not null"`
	CredentialID    []byte    `gorm:"type:bytea;not null;unique"`
	PublicKey       []byte    `gorm:"type:bytea;not null"`
	SignCount       int64     `gorm:"type:bigint;not null"`
	AttestationType string    `gorm:"type:text"`
	Transports      []string  `gorm:"type:text[]"`
	CreatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

type Session struct {
	ID         uint           `gorm:"primaryKey"`
	SessionKey string         `gorm:"type:text;not null;unique"`
	UserID     uint           `gorm:"not null"`
	Data       datatypes.JSON `gorm:"type:jsonb;not null"`
	CreatedAt  time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	ExpiresAt  time.Time      `gorm:"not null"`
}
