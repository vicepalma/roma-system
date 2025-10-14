// internal/security/invite.go
package security

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type InviteClaims struct {
	CoachID       string `json:"coach_id"`
	DiscipleEmail string `json:"disciple_email"`
	jwt.RegisteredClaims
}

// Emite un JWT de invitación (audience=invite, typ=invite)
func SignInvite(coachID, discipleEmail string, ttl time.Duration) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET missing")
	}
	now := time.Now()
	claims := InviteClaims{
		CoachID:       coachID,
		DiscipleEmail: discipleEmail,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   coachID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			Audience:  jwt.ClaimStrings{"invite"},
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tok.Header["typ"] = "invite"
	return tok.SignedString([]byte(secret))
}

// Valida y parsea el código de invitación
func ParseInvite(code string) (*InviteClaims, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET missing")
	}
	var claims InviteClaims
	tok, err := jwt.ParseWithClaims(
		code,
		&claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithAudience("invite"), // <- valida aud automáticamente en v5
	)
	if err != nil {
		return nil, err
	}
	if !tok.Valid {
		return nil, errors.New("invalid invite code")
	}
	return &claims, nil
}
