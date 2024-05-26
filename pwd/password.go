package pwd

import (
	cryptorand "crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"math/rand"

	"strings"

	"golang.org/x/crypto/pbkdf2"
)

const (
	saltLen    = 16
	iterations = 100
	keyLen     = 32
	alpha      = "pbkdf2-sha512"
)

func GenPassword(pass string) string {
	hashFunction := sha512.New
	salt := generateSalt(saltLen)
	encodedPwd := hex.EncodeToString(pbkdf2.Key([]byte(pass), salt, iterations, keyLen, hashFunction))
	newEncodedPwd := fmt.Sprintf("$%s$%s$%s", alpha, salt, encodedPwd)
	return newEncodedPwd
}

func VerifyPassword(rePass, encodedPwd string) bool {
	hashFunction := sha512.New
	infos := strings.Split(encodedPwd, "$")
	salt := []byte(infos[2])
	newPass := hex.EncodeToString(pbkdf2.Key([]byte(rePass), salt, iterations, keyLen, hashFunction))
	if infos[3] == newPass {
		return true
	}
	return false
}

func generateSalt(length int) []byte {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	salt := make([]byte, length)
	if _, err := cryptorand.Read(salt); err != nil {
		rand.Read(salt)
	}
	for key, val := range salt {
		salt[key] = alphanum[val%byte(len(alphanum))]
	}
	return salt
}
