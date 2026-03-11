package auth

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type jwksResponse struct {
	Keys []jwkEntry `json:"keys"`
}

type jwkEntry struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Crv string `json:"crv"`
	X   string `json:"x"`
	Alg string `json:"alg"`
	Use string `json:"use"`
}

// FetchJWKS fetches the JWKS from the identity service and returns a map of
// kid -> Ed25519 public key. Fails fast if the identity service is unreachable
// or returns unexpected key types.
func FetchJWKS(identityURL string) (map[string]ed25519.PublicKey, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(identityURL + "/v1/.well-known/jwks.json")
	if err != nil {
		return nil, fmt.Errorf("fetch JWKS from %s: %w", identityURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch JWKS: unexpected status %d", resp.StatusCode)
	}

	var jwks jwksResponse
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, fmt.Errorf("decode JWKS: %w", err)
	}

	keys := make(map[string]ed25519.PublicKey)
	for _, entry := range jwks.Keys {
		if entry.Kty != "OKP" || entry.Crv != "Ed25519" {
			continue
		}
		if entry.Kid == "" || entry.X == "" {
			continue
		}

		pubBytes, err := base64.RawURLEncoding.DecodeString(entry.X)
		if err != nil {
			return nil, fmt.Errorf("decode public key for kid %s: %w", entry.Kid, err)
		}
		if len(pubBytes) != ed25519.PublicKeySize {
			return nil, fmt.Errorf("invalid public key size for kid %s: got %d, want %d", entry.Kid, len(pubBytes), ed25519.PublicKeySize)
		}

		keys[entry.Kid] = ed25519.PublicKey(pubBytes)
	}

	if len(keys) == 0 {
		return nil, fmt.Errorf("no Ed25519 keys found in JWKS from %s", identityURL)
	}

	return keys, nil
}
