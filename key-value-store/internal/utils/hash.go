package utils

import "crypto/sha256"

func HashKey(key string) uint32 {
    hash := sha256.Sum256([]byte(key))
    return uint32(hash[0]) | uint32(hash[1])<<8 | uint32(hash[2])<<16 | uint32(hash[3])<<24
}