package tool

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// Encrypt 使用AES算法对数据进行加密
//
// 该函数使用AES算法的CFB模式对明文进行加密。
// 加密过程包括：
// 1. 生成随机的初始化向量(IV)
// 2. 使用CFB模式加密数据
// 3. 将IV拼接到密文头部
// 4. 对结果进行Base64 URL编码
//
// 参数:
//   - plaintext: 需要加密的明文数据
//   - key: 加密密钥，长度必须为16、24或32字节（分别对应AES-128, AES-192, AES-256）
//
// 返回值:
//   - string: Base64 URL编码后的加密字符串（包含IV）
//   - error: 如果加密过程中发生错误（如key长度不合法、随机数生成失败），返回相应的错误信息
//
// 使用示例:
//
//	key := []byte("0123456789abcdef0123456789abcdef") // 32 bytes for AES-256
//	plaintext := []byte("Hello, World!")
//
//	encrypted, err := Encrypt(plaintext, key)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("Encrypted:", encrypted)
//
// 注意事项:
//   - 密钥长度必须严格符合AES标准（16/24/32字节）
//   - 每次加密生成的密文都不同（因为使用了随机IV）
//   - 返回的字符串是URL安全的Base64编码
func Encrypt(plaintext, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// Return the encrypted data as a base64 encoded string
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 对AES加密的数据进行解密
//
// 该函数用于解密由Encrypt函数生成的密文。
// 解密过程包括：
// 1. 对输入字符串进行Base64 URL解码
// 2. 从解码后的数据头部提取IV
// 3. 使用AES CFB模式解密剩余数据
//
// 参数:
//   - ciphertext: Base64 URL编码的密文字符串
//   - key: 解密密钥，必须与加密时使用的密钥完全一致
//
// 返回值:
//   - string: 解密后的明文字符串
//   - error: 如果解密失败（如Base64解码失败、密文太短、Key无效），返回错误信息
//
// 使用示例:
//
//	key := []byte("0123456789abcdef0123456789abcdef")
//	encrypted := "..." // Encrypt函数的输出
//
//	decrypted, err := Decrypt(encrypted, key)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("Decrypted:", decrypted)
//
// 注意事项:
//   - 如果密文被篡改或截断，解密可能会失败或产生乱码
//   - 密钥必须与加密时使用的完全一致
func Decrypt(ciphertext string, key []byte) (string, error) {
	ciphertextBytes, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertextBytes) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := ciphertextBytes[:aes.BlockSize]
	ciphertextBytes = ciphertextBytes[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertextBytes, ciphertextBytes)

	return string(ciphertextBytes), nil
}