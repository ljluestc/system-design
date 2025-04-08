package utils

import (
    "crypto/sha256"
    "encoding/hex"
)

func HashToInt(s string) int {
    hash := sha256.Sum256([]byte(s))
    hashStr := hex.EncodeToString(hash[:])
    return int(hashStr[0])
}