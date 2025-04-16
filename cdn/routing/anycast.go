package routing

import (
    "crypto/md5"
    "strconv"
)

func GetEdgeID(key string, numEdges int) int {
    hash := md5.Sum([]byte(key))
    hashInt, _ := strconv.ParseInt(string(hash[:8]), 16, 64)
    return int(hashInt) % numEdges
}