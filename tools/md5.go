package tools

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5 returns the md5 hash of a []byte
func MD5(input []byte) string {
	hasher := md5.New()
	hasher.Write(input)
	return hex.EncodeToString(hasher.Sum(nil))
}
