package crypto

import "C"
import (
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/chacha20poly1305"
)

const (
	secret = "EduHacks2020.*"
)

var key = sha256.Sum256([]byte(secret))

// 包含一个 AEAD 加密实例
type ChaCha20Poly1305 struct {
	aead cipher.AEAD
}

// Init 初始化密匙
func (c *ChaCha20Poly1305) Init() {
	c.aead, _ = chacha20poly1305.New(key[:])
}

// EncryptedToBase64 加密到 Base64 输出
func (c *ChaCha20Poly1305) EncryptedToBase64(msg string) string {
	nonce := make([]byte, c.aead.NonceSize(), c.aead.NonceSize()+len(msg)+c.aead.Overhead())
	if _, err := rand.Read(nonce); err != nil {
		log.Error(err.Error())
	}
	ciphertext := c.aead.Seal(nonce, nonce, []byte(msg), nil)
	return base64.URLEncoding.EncodeToString(ciphertext)
}

// DecryptedFromBase64 从 Base64 解密
func (c *ChaCha20Poly1305) DecryptedFromBase64(text string) ([]byte, error) {
	// 先对其进行 base64 解码
	decodeBytes, err := base64.URLEncoding.DecodeString(text)
	if err != nil {
		log.Error(err.Error())
	}
	if len(decodeBytes) < c.aead.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}
	nonce, ciphertext := decodeBytes[:c.aead.NonceSize()], decodeBytes[c.aead.NonceSize():]
	plaintext, err := c.aead.Open(nil, nonce, ciphertext, nil)
	return plaintext, err
}

// EncryptedToHex 加密到十六进制输出
func (c *ChaCha20Poly1305) EncryptedToHex(msg string) string {
	nonce := make([]byte, c.aead.NonceSize(), c.aead.NonceSize()+len(msg)+c.aead.Overhead())
	if _, err := rand.Read(nonce); err != nil {
		log.Error(err.Error())
	}
	ciphertext := c.aead.Seal(nonce, nonce, []byte(msg), nil)
	return hex.EncodeToString(ciphertext)
}

// DecryptedFromHex 从十六进制解密
func (c *ChaCha20Poly1305) DecryptedFromHex(text string) ([]byte, error) {
	// 先对十六禁止进行解码
	decodeBytes, err := hex.DecodeString(text)
	if err != nil {
		log.Error(err.Error())
	}
	if len(decodeBytes) < c.aead.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}
	nonce, ciphertext := decodeBytes[:c.aead.NonceSize()], decodeBytes[c.aead.NonceSize():]
	plaintext, err := c.aead.Open(nil, nonce, ciphertext, nil)
	return plaintext, err
}
