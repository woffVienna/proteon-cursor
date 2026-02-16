package auth

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTIssuer issues Ed25519-signed JWTs.
type JWTIssuer struct {
	kid        string
	publicKey  ed25519.PublicKey
	privateKey ed25519.PrivateKey
	issuer     string
	audience   string
}

// NewJWTIssuer creates an Ed25519 JWT issuer.
// For production, load keys from env/Secrets Manager/KMS.
func NewJWTIssuer() (*JWTIssuer, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return &JWTIssuer{
		kid:        "dev-identity-001",
		publicKey:  pub,
		privateKey: priv,
		issuer:     IssuerFromEnv(),
		audience:   AudienceFromEnv(),
	}, nil
}

// Issue implements interfaces.TokenIssuer.
func (j *JWTIssuer) Issue(_ context.Context, userID, tenant string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"iss":    j.issuer,
		"aud":    j.audience,
		"sub":    userID,
		"iat":    now.Unix(),
		"nbf":    now.Unix(),
		"exp":    now.Add(ttl).Unix(),
		"tenant": tenant,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = j.kid
	return token.SignedString(j.privateKey)
}

// PublicKey implements interfaces.TokenIssuer.
func (j *JWTIssuer) PublicKey() ed25519.PublicKey {
	return j.publicKey
}

// Kid implements interfaces.TokenIssuer.
func (j *JWTIssuer) Kid() string {
	return j.kid
}
