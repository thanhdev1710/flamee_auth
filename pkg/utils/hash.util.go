package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"

	"github.com/thanhdev1710/flamee_auth/global"
)

// Hàm mã hóa dữ liệu bằng AES-GCM
func Encrypt(data string) (string, error) {
	// Tạo block cipher
	block, err := aes.NewCipher([]byte(global.Config.Email.Secret))
	if err != nil {
		return "", err
	}

	// Tạo nonce (IV)
	nonce := make([]byte, 12) // GCM yêu cầu nonce có độ dài 12 byte
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Tạo AES-GCM cipher
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Mã hóa dữ liệu
	ciphertext := aesGCM.Seal(nil, nonce, []byte(data), nil)

	// Trả về mã hóa dưới dạng hex
	return hex.EncodeToString(append(nonce, ciphertext...)), nil
}

// Hàm giải mã dữ liệu bằng AES-GCM
func Decrypt(ciphertextHex string) (string, error) {
	// Giải mã từ hex về dạng byte
	ciphertext, err := hex.DecodeString(ciphertextHex)
	if err != nil {
		return "", err
	}

	// Tách nonce và dữ liệu mã hóa
	nonce, ciphertext := ciphertext[:12], ciphertext[12:]

	// Tạo block cipher
	block, err := aes.NewCipher([]byte(global.Config.Email.Secret))
	if err != nil {
		return "", err
	}

	// Tạo AES-GCM cipher
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Giải mã dữ liệu
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errors.New("failed to decrypt data")
	}

	return string(plaintext), nil
}
