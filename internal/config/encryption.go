package config

import (
	"os"
	"time"
)

var(
	JWT_SIGNATURE_KEY	string	
	ENCRYPTION_KEY	string	
)

const (
	// AccessTokenExpiry  = time.Minute * 15   // 15 minutes
	AccessTokenExpiry  = time.Second * 15   // 15 second
	RefreshTokenExpiry = time.Hour * 24 * 7 // 7 days
)

func InitEncryptionVars() {
	JWT_SIGNATURE_KEY = os.Getenv("JWT_SIGNATURE_KEY"+Prefix)
	ENCRYPTION_KEY = os.Getenv("ENCRYPTION_KEY"+Prefix)
}