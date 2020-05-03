package crypto_utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

func GetPasswordHash(password string, salt string) string {
	password = fmt.Sprint(password, salt)
	hash := sha256.New()
	defer hash.Reset()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

func GetSha256(input string) string {
	hash := sha256.New()
	defer hash.Reset()
	hash.Write([]byte(input))
	return hex.EncodeToString(hash.Sum(nil))
}

func SaltText() string {
	salt := GetSha256(fmt.Sprintf("%d", time.Now().UnixNano()))
	return salt
}
