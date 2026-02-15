package jwtverifier

import (
	"crypto/ed25519"
	"errors"
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

var (
	ErrUnauthorized     = errors.New("unauthorized")
	ErrUnsupportedAlg   = errors.New("unsupported jwt alg")
	ErrUnknownKeyID     = errors.New("unknown kid")
	ErrMissingKeyID     = errors.New("missing kid")
	ErrInvalidToken     = errors.New("invalid token")
	ErrInvalidClaims    = errors.New("invalid claims")
	ErrMissingSubject   = errors.New("missing sub")
	ErrMissingTenant    = errors.New("missing tenant")
	ErrAudienceMismatch = errors.New("audience mismatch")
	ErrIssuerMismatch   = errors.New("issuer mismatch")
	ErrTokenExpired     = errors.New("token expired")
	ErrTokenNotYetValid = errors.New("token not yet valid")
)

type Config struct {
	// Allowed issuer/audience. If empty, issuer/audience checks are skipped.
	Issuer   string
	Audience string

	// Map of kid -> Ed25519 public key.
	// This matches how JWT header "kid" selects the verification key.
	Keys map[string]ed25519.PublicKey

	// Optional clock skew to tolerate (e.g. 30s) when validating time claims.
	Leeway time.Duration
}

// Claims is the minimum set of claims you care about platform-wide.
// Keep it small; treat everything else as optional.
type Claims struct {
	Subject   string
	Tenant    string
	Scopes    []string
	SessionID string // optional; empty if not present
	KeyID     string
	ExpiresAt time.Time
	IssuedAt  time.Time
}

type Verifier struct {
	cfg Config
}

func New(cfg Config) *Verifier {
	return &Verifier{cfg: cfg}
}

// Verify verifies signature + standard claims and extracts your domain claims.
// It returns ErrUnauthorized (wrapped) for any auth failure.
func (v *Verifier) Verify(rawToken string) (Claims, error) {
	if rawToken == "" {
		return Claims{}, wrap(ErrUnauthorized, ErrInvalidToken)
	}

	// Parse headers first to get kid. jwt.ParseWithClaims can expose header via token.Header
	// but we want to enforce presence early.
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		// Enforce algorithm
		if t.Method == nil || t.Method.Alg() != jwt.SigningMethodEdDSA.Alg() {
			return nil, wrap(ErrUnauthorized, ErrUnsupportedAlg)
		}

		kidAny, ok := t.Header["kid"]
		if !ok {
			return nil, wrap(ErrUnauthorized, ErrMissingKeyID)
		}
		kid, ok := kidAny.(string)
		if !ok || kid == "" {
			return nil, wrap(ErrUnauthorized, ErrMissingKeyID)
		}

		pub, ok := v.cfg.Keys[kid]
		if !ok {
			return nil, wrap(ErrUnauthorized, ErrUnknownKeyID)
		}

		return pub, nil
	}

	parserOpts := []jwt.ParserOption{
		jwt.WithValidMethods([]string{jwt.SigningMethodEdDSA.Alg()}),
	}
	if v.cfg.Leeway > 0 {
		parserOpts = append(parserOpts, jwt.WithLeeway(v.cfg.Leeway))
	}
	if v.cfg.Issuer != "" {
		parserOpts = append(parserOpts, jwt.WithIssuer(v.cfg.Issuer))
	}
	if v.cfg.Audience != "" {
		parserOpts = append(parserOpts, jwt.WithAudience(v.cfg.Audience))
	}

	parser := jwt.NewParser(parserOpts...)

	claims := jwt.MapClaims{}
	tok, err := parser.ParseWithClaims(rawToken, claims, keyFunc)
	if err != nil {
		// Map common time errors (best-effort).
		if errors.Is(err, jwt.ErrTokenExpired) {
			return Claims{}, wrap(ErrUnauthorized, ErrTokenExpired)
		}
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return Claims{}, wrap(ErrUnauthorized, ErrTokenNotYetValid)
		}
		return Claims{}, wrap(ErrUnauthorized, ErrInvalidToken)
	}
	if tok == nil || !tok.Valid {
		return Claims{}, wrap(ErrUnauthorized, ErrInvalidToken)
	}

	// kid
	kid, _ := tok.Header["kid"].(string)

	// sub
	sub, _ := claims["sub"].(string)
	if sub == "" {
		return Claims{}, wrap(ErrUnauthorized, ErrMissingSubject)
	}

	// tenant (your custom claim)
	tenant, _ := claims["tenant"].(string)
	if tenant == "" {
		return Claims{}, wrap(ErrUnauthorized, ErrMissingTenant)
	}

	// scopes (optional): support "scope": "a b" or "scopes": ["a","b"]
	var scopes []string
	if scopeStr, ok := claims["scope"].(string); ok && scopeStr != "" {
		scopes = splitScopes(scopeStr)
	} else if scopeAny, ok := claims["scopes"]; ok {
		if arr, ok := scopeAny.([]interface{}); ok {
			for _, v := range arr {
				if s, ok := v.(string); ok && s != "" {
					scopes = append(scopes, s)
				}
			}
		}
	}

	// session id (optional): "sid"
	sid, _ := claims["sid"].(string)

	// exp/iat (optional but useful)
	var exp time.Time
	if expF, ok := claims["exp"].(float64); ok && expF > 0 {
		exp = time.Unix(int64(expF), 0)
	}
	var iat time.Time
	if iatF, ok := claims["iat"].(float64); ok && iatF > 0 {
		iat = time.Unix(int64(iatF), 0)
	}

	return Claims{
		Subject:   sub,
		Tenant:    tenant,
		Scopes:    scopes,
		SessionID: sid,
		KeyID:     kid,
		ExpiresAt: exp,
		IssuedAt:  iat,
	}, nil
}

func splitScopes(s string) []string {
	out := make([]string, 0, 8)
	start := -1
	for i := 0; i < len(s); i++ {
		if s[i] != ' ' {
			if start == -1 {
				start = i
			}
			continue
		}
		if start != -1 {
			out = append(out, s[start:i])
			start = -1
		}
	}
	if start != -1 {
		out = append(out, s[start:])
	}
	return out
}

func wrap(top, cause error) error {
	return fmt.Errorf("%w: %v", top, cause)
}
