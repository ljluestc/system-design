package encoder

import (
    "math/big"
)

// Base58Alphabet defines the characters used
const Base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// ToBase58 converts a uint64 to base-58
func ToBase58(num uint64) string {
    if num == 0 {
        return string(Base58Alphabet[0])
    }
    var result []byte
    base := big.NewInt(58)
    zero := big.NewInt(0)
    n := big.NewInt(int64(num))
    for n.Cmp(zero) > 0 {
        rem := new(big.Int)
        n.DivMod(n, base, rem)
        result = append([]byte{Base58Alphabet[rem.Int64()]}, result...)
    }
    return string(result)
}