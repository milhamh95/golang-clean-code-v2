package md5

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/pkg/errors"
)

// Generate computes the MD5 checksum of the given string
func Generate(s string) (string, error) {
	hasher := md5.New()
	_, err := hasher.Write([]byte(s))
	if err != nil {
		return "", errors.Wrap(err, "error generating md5")
	}
	result := hex.EncodeToString(hasher.Sum(nil))
	return result, nil
}
