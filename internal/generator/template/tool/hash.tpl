package tool

import (
	"crypto/md5"
	"encoding/hex"
	"hash/fnv"
)

// MD5 计算MD5
func MD5(data string) string {
	// 计算 MD5
	hash := md5.Sum([]byte(data))
	// 转成 16 进制字符串
	hashStr := hex.EncodeToString(hash[:])
	return hashStr
}

// HashToInt64 将字符串 hash 成 int64
// 使用 FNV-1a 算法，性能好且分布均匀
func HashToInt64(data string) int64 {
	h := fnv.New64a()
	_, err := h.Write([]byte(data))
	if err != nil {
		return 0
	}
	return int64(h.Sum64())
}
