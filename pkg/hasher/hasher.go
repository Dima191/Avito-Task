package hasher

import (
	"golang.org/x/crypto/bcrypt"
	"unsafe"
)

// Hash returns the bcrypt hash
func Hash(data string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(unsafe.Slice(unsafe.StringData(data), len(data)), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return unsafe.String(unsafe.SliceData(hash), len(hash)), nil
}

// Compare compares a bcrypt hashed password with its possible
// plaintext equivalent. Returns nil on success, or an error on failure.
func Compare(data, hash string) error {
	passwordBytes := unsafe.Slice(unsafe.StringData(data), len(data))
	hashBytes := unsafe.Slice(unsafe.StringData(hash), len(hash))
	if err := bcrypt.CompareHashAndPassword(hashBytes, passwordBytes); err != nil {
		return err
	}

	return nil
}
