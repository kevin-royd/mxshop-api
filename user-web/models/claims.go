package models

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type CustomClaims struct {
	ID          uint   `json:"id"`
	NickName    string `json:"nickname"`
	AuthorityId uint   `json:"authority_id"`
	jwt.RegisteredClaims
}

// Valid Implement the Valid method for CustomClaims
func (c CustomClaims) Valid() error {
	if c.ExpiresAt != nil && c.ExpiresAt.Before(time.Now()) {
		return jwt.ErrTokenExpired
	}
	return nil
}
