package entity

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"

	AccessTokenLifeTime  = 30 // minutes
	RefreshTokenLifeTime = 48 // hours

	ClaimKeyTokenId       = "token_id"
	ClaimKeyPairedTokenId = "paired_token_id"
	ClaimKeyUserId        = "user_id"
	ClaimKeyType          = "type"
	ClaimKeyExpiredAt     = "expired_at"

	ContextTokenIDKey       = "token_id"
	ContextPairedTokenIdKey = "paired_token_id"
	ContextTokenTypeKey     = "token_type"
)

type Token struct {
	ID          uuid.UUID
	Value       string
	ExpiredTime time.Time
}

type TokenPair struct {
	AccessToken  Token
	RefreshToken Token
}

type TokenPairFormat struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (tp *TokenPair) Format() interface{} {
	return &TokenPairFormat{
		AccessToken:  tp.AccessToken.Value,
		RefreshToken: tp.RefreshToken.Value,
	}
}
