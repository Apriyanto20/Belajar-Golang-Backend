package main

import (
	"math/rand"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func main() {
	r := gin.Default()

	// ======= ENABLE CORS =======
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Boleh diakses dari mana saja
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// ROUTE TUNGGAL
	r.Any("/", handleAll)

	r.Run(":8080")
}

// =========================
// HANDLE ALL (GENERATE KEY / ENCRYPT / DECRYPT)
// =========================
func handleAll(c *gin.Context) {
	// Jika GET → generate key
	if c.Request.Method == "GET" {
		c.JSON(http.StatusOK, gin.H{"key": generateRandomKey()})
		return
	}

	// Jika POST → encrypt / decrypt
	if c.Request.Method == "POST" {
		var body map[string]string
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		// Kalau body kosong → generate key juga bisa
		if len(body) == 0 {
			c.JSON(http.StatusOK, gin.H{"key": generateRandomKey()})
			return
		}

		key, ok := body["key"]
		if !ok || len(key) != len(alphabet) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing key"})
			return
		}

		// ENCRYPT
		if plainText, ok := body["plain_text"]; ok {
			encrypted := processText(plainText, alphabet, key)
			c.JSON(http.StatusOK, gin.H{"encrypted_text": encrypted})
			return
		}

		// DECRYPT
		if encryptedText, ok := body["encrypted_text"]; ok {
			decrypted := processText(encryptedText, key, alphabet)
			c.JSON(http.StatusOK, gin.H{"plain_text": decrypted})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": "body must contain plain_text or encrypted_text"})
		return
	}

	// OPTIONS → untuk CORS preflight
	c.Status(http.StatusOK)
}

// =========================
// GENERATE RANDOM KEY
// =========================
func generateRandomKey() string {
	rand.Seed(time.Now().UnixNano())
	letters := strings.Split(alphabet, "")
	rand.Shuffle(len(letters), func(i, j int) {
		letters[i], letters[j] = letters[j], letters[i]
	})
	return strings.Join(letters, "")
}

// =========================
// PROCESS TEXT (ENCRYPT / DECRYPT)
// =========================
func processText(input, from, to string) string {
	var result strings.Builder
	for _, char := range input {
		if char == ' ' {
			result.WriteRune(char)
			continue
		}
		isUpper := unicode.IsUpper(char)
		lowerChar := unicode.ToLower(char)
		index := strings.IndexRune(from, lowerChar)
		if index == -1 {
			result.WriteRune(char)
			continue
		}
		newChar := rune(to[index])
		if isUpper {
			newChar = unicode.ToUpper(newChar)
		}
		result.WriteRune(newChar)
	}
	return result.String()
}
