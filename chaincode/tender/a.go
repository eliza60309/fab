package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

func main() {
	fmt.Println("Enter mode:")
	var mode int
	fmt.Scanln(&mode)
	if(mode == 0) {
		fmt.Println("Enter input:")
		var originalText string
		fmt.Scanln(&originalText)
		fmt.Println("Enter key:")
		var keyid int
		fmt.Scanln(&keyid)
		keytext := ""
		if keyid == 1 {
			keytext = "NRxPLLBMIJ"
		} else if keyid == 2 {
			keytext = "jx52X9ASHP"
		} else if keyid == 3 {
			keytext = "posARbsxeV"
		} else {
			return 
		}
		i := len(keytext)
		fmt.Println("key size = ", i)
		key := []byte(keytext)[:16]
		cryptoText := encrypt(key, originalText)
		fmt.Println(cryptoText)
	} else {
		fmt.Println("Enter input:")
		var originalText string
		fmt.Scanln(&originalText)
		fmt.Println("Enter key:")
		var keyid int
		fmt.Scanln(&keyid)
		keytext := ""
		if keyid == 1 {
			keytext = "NRxPLLBMIJ"
		} else if keyid == 2 {
			keytext = "jx52X9ASHP"
		} else if keyid == 3 {
			keytext = "posARbsxeV"
		} else {
			return 
		}
		key := []byte(keytext)[:16]
		cryptoText := decrypt(key, originalText)
		fmt.Println(cryptoText)

	}
	//originalText := "encrypt this golang"

	// encrypt value to base64
	//cryptoText := encrypt(key, originalText)
	//fmt.Println(cryptoText)

	// encrypt base64 crypto to original value
	//text := decrypt(key, cryptoText)
	//fmt.Printf(text)
}

// encrypt string to base64 crypto using AES
func encrypt(key []byte, text string) string {
	// key := []byte(keyText)
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext)
}

// decrypt from base64 to decrypted string
func decrypt(key []byte, cryptoText string) string {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}
