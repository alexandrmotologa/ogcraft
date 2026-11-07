package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"sort"
	"strings"
)

// GenerateSignature generates an HMAC-SHA256 signature for URL parameters.
func GenerateSignature(secret, template string, params map[string]string) string {
	if secret == "" {
		return ""
	}

	var keys []string
	for k := range params {
		if k == "sig" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	sb.WriteString(template)
	sb.WriteString("?")
	for _, k := range keys {
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(params[k])
		sb.WriteString("&")
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(sb.String()))
	return hex.EncodeToString(mac.Sum(nil))
}

// VerifySignature validates that the provided signature matches the HMAC-SHA256 hash.
func VerifySignature(secret, template string, params map[string]string, sig string) bool {
	if secret == "" {
		return true // Signing not required
	}
	if sig == "" {
		return false
	}

	expectedSig := GenerateSignature(secret, template, params)
	return subtle.ConstantTimeCompare([]byte(expectedSig), []byte(sig)) == 1
}
