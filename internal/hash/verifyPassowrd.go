package hash

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

func VerifyPassword(password, encodedHash string) (bool, error) {
	// 1. Split the stored string to get parameters and salt
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, fmt.Errorf("invalid hash format")
	}

	var m, t, p uint32
	fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &m, &t, &p)

	salt, _ := base64.RawStdEncoding.DecodeString(parts[4])
	decodedHash, _ := base64.RawStdEncoding.DecodeString(parts[5])

	// 2. Re-hash the login attempt with the SAME salt and params
	comparisonHash := argon2.IDKey([]byte(password), salt, t, m, uint8(p), uint32(len(decodedHash)))

	// 3. Constant-time comparison to prevent timing attacks
	return subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1, nil
}
