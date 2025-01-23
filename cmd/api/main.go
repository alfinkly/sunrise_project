package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sunrise_project/internal/platform"
	"sunrise_project/internal/repository"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file, using environment variables instead")
	}

	db, err := platform.NewPostgresDB()
	if err != nil {
		log.Fatalf("Failed to connect: %s", err)
	}
	err = db.AutoMigrate(&repository.Location{})
	if err != nil {
		log.Fatalf("Failed to Migrate Database: %s", err)
	}

	key := []byte("dog-alfinkly-bird-apple-taxa-dot")

	encodedCiphertext, err := ioutil.ReadFile("secret_data.txt")
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	ciphertext, err := base64.StdEncoding.DecodeString(string(encodedCiphertext))
	if err != nil {
		log.Fatalf("Failed to decode base64: %v", err)
	}

	decryptedText, err := decrypt(ciphertext, key)
	if err != nil {
		log.Fatalf("Failed to decrypt: %v", err)
	}

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, decryptedText)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port - %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}

func decrypt(ciphertext []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
