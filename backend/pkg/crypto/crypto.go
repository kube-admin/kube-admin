// Package crypto 提供对称加密能力，用于保护存储在数据库中的集群敏感凭据。
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

var (
	encryptKey  []byte
	errNotInit  = errors.New("crypto 包未初始化，请先调用 crypto.Init")
	errEmptyKey = errors.New("加密密钥不能为空")
)

// Init 初始化加密密钥。传入任意长度密钥，内部用 SHA-256 派生 32 字节
// 密钥用于 AES-256-GCM。重复调用以最后一次为准。
func Init(key string) error {
	if key == "" {
		return errEmptyKey
	}
	sum := sha256.Sum256([]byte(key))
	encryptKey = sum[:]
	return nil
}

// Encrypt 加密明文，返回 base64 编码的密文（含 nonce）。
// 空字符串原样返回，避免空值被加密后产生歧义。
func Encrypt(plain string) (string, error) {
	if plain == "" {
		return "", nil
	}
	if len(encryptKey) == 0 {
		return "", errNotInit
	}

	block, err := aes.NewCipher(encryptKey)
	if err != nil {
		return "", fmt.Errorf("创建 AES cipher 失败: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("创建 GCM 失败: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("生成 nonce 失败: %w", err)
	}

	// nonce 前置，解密时按 NonceSize 截取
	ciphertext := gcm.Seal(nonce, nonce, []byte(plain), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 解密 Encrypt 产出的 base64 密文，返回明文。
// 空字符串原样返回。无法识别的旧明文（未加密）会返回错误，调用方可降级处理。
func Decrypt(encoded string) (string, error) {
	if encoded == "" {
		return "", nil
	}
	if len(encryptKey) == 0 {
		return "", errNotInit
	}

	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("base64 解码失败: %w", err)
	}

	block, err := aes.NewCipher(encryptKey)
	if err != nil {
		return "", fmt.Errorf("创建 AES cipher 失败: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("创建 GCM 失败: %w", err)
	}

	if len(data) < gcm.NonceSize() {
		return "", errors.New("密文长度不足")
	}

	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("解密失败: %w", err)
	}
	return string(plain), nil
}
