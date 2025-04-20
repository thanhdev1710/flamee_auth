package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"

	"github.com/thanhdev1710/flamee_auth/global"
	"golang.org/x/crypto/bcrypt"
)

// Hàm mã hóa dữ liệu bằng AES-GCM
func Encrypt(data string) (string, error) {
	// Tạo block cipher
	block, err := aes.NewCipher([]byte(global.Config.Email.Secret))
	if err != nil {
		return "", errors.New("không thể tạo cipher từ khóa bí mật")
	}

	// Tạo nonce (IV)
	nonce := make([]byte, 12) // GCM yêu cầu nonce có độ dài 12 byte
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", errors.New("không thể tạo nonce ngẫu nhiên")
	}

	// Tạo AES-GCM cipher
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.New("không thể khởi tạo AES-GCM")
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
		return "", errors.New("chuỗi mã hóa không hợp lệ")
	}

	// Kiểm tra độ dài trước khi tách nonce
	if len(ciphertext) < 12 {
		return "", errors.New("dữ liệu mã hóa không đầy đủ")
	}

	// Tách nonce và dữ liệu mã hóa
	nonce, ciphertext := ciphertext[:12], ciphertext[12:]

	// Tạo block cipher
	block, err := aes.NewCipher([]byte(global.Config.Email.Secret))
	if err != nil {
		return "", errors.New("không thể tạo cipher từ khóa bí mật")
	}

	// Tạo AES-GCM cipher
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.New("không thể khởi tạo AES-GCM")
	}

	// Giải mã dữ liệu
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errors.New("giải mã dữ liệu không thành công")
	}

	return string(plaintext), nil
}

func CompareHashAndPassword(hashedPassword, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return errors.New("mật khẩu không đúng")
	}
	return nil
}

func GenerateFromPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func HexString(str string) string {
	return hex.EncodeToString([]byte(str))
}
