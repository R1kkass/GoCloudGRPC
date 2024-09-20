package helpers

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"math/big"
)

func GenerateSecretKey(p *big.Int, b string, A string) (*big.Int, error) {
	B, _ := new(big.Int).SetString(b, 0)
	secretKey := new(big.Int)
	userKey := new(big.Int)
	userKey, _ = userKey.SetString(A, 0)
	secretKey = secretKey.Exp(userKey, B, nil)
	secretKey = secretKey.Mod(secretKey, p)
	return secretKey, nil
}



func generateIV2(cryptoText string) string {
	decodedData, _ := base64.RawURLEncoding.DecodeString(cryptoText)
	ivBytes := decodedData[:16]
	return base64.StdEncoding.EncodeToString(ivBytes)
}

func Crypt(isEncrypt bool, data []byte, secretKey string) []byte {
	iv := generateIV2(secretKey)
	ivBytes, _ := base64.StdEncoding.DecodeString(iv)
	key := []byte(secretKey)

	block, _ := aes.NewCipher(key)
	mode := cipher.NewCBCEncrypter(block, ivBytes)

	if !isEncrypt {
		mode = cipher.NewCBCDecrypter(block, ivBytes)
	}

	paddedData := pad(data)
	if isEncrypt {
		ciphertext := make([]byte, len(paddedData))
		mode.CryptBlocks(ciphertext, paddedData)
		return ciphertext
	} else {
		plaintext := make([]byte, len(data))
		mode.CryptBlocks(plaintext, data)
		return unpad(plaintext)
	}
}

func Decrypt(data string, secretKey string) string{
	text, _ := base64.StdEncoding.DecodeString(data)
	decrypted := Crypt(false, text, secretKey)
	return string(decrypted)
}

func Encrypt(data string, secretKey string) string{
	plaintext := []byte(data)
	chipText := Crypt(true, plaintext,secretKey)
	return base64.StdEncoding.EncodeToString(chipText)
}

func pad(data []byte) []byte {
	padding := aes.BlockSize - len(data)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

func unpad(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}
