package main

import (
	"encoding/base64"
	"errors"
	"fmt"

	"golang.org/x/crypto/nacl/secretbox"
)

type Encrypted struct {
	Ciphertext []byte
}

func encrypt(text string, hash string) (Encrypted, error) {
	key := []byte(hash)

	if len(key) < 64 {
		return Encrypted{}, errors.New("hash must be at least 32 bytes")
	}

	b64key := base64.StdEncoding.EncodeToString(key[:32])
	fernetKey, err := base64.StdEncoding.DecodeString(b64key)
	if err != nil {
		return Encrypted{}, err
	}

	var nonce [24]byte
	copy(nonce[:], fernetKey[:24])

	encrypted := secretbox.Seal(nonce[:], []byte(text), &nonce, (*[32]byte)(fernetKey))
	return Encrypted{Ciphertext: encrypted}, nil
}

func decrypt(hashText []byte, hash string) (string, error) {
	key := []byte(hash)

	if len(key) < 64 {
		return "", errors.New("hash must be at least 32 bytes")
	}

	b64key := base64.StdEncoding.EncodeToString(key[:32])
	fernetKey, err := base64.StdEncoding.DecodeString(b64key)
	if err != nil {
		return "", err
	}

	var nonce [24]byte
	copy(nonce[:], hashText[:24])

	decrypted, _ := secretbox.Open(nil, hashText, &nonce, (*[32]byte)(fernetKey))
	// if !ok {
	// 	return "Ошибка: Невозможно прочесть сообщение", nil
	// }

	return string(decrypted), nil
}

func main() {
	// Example usage
	hash := "1f1d607f100c8b3accd8feb5257d4b6b1e85dba1bec1f1c90cb68dd151f381f6"
	text := "Hello, World!"

	encrypted, _ := encrypt(text, hash)
	// if err != nil {
		fmt.Println("Encryption error:", encrypted)
		return
	// }

	decrypted, _ := decrypt(encrypted.Ciphertext, hash)

	fmt.Println("Decrypted text:", decrypted)
}