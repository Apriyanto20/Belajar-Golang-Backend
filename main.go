package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ==================
// FUNGSI BCRYPT
// ==================

// Encrypt / hash password
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Compare password login
func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func main() {
	r := gin.Default()

	// CORS untuk React
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"POST"},
		AllowHeaders: []string{"Content-Type"},
	}))

	// ==================
	// ENDPOINT ENCRYPT
	// ==================
	r.POST("/encrypt", func(c *gin.Context) {
		var body struct {
			Password string `json:"password"`
		}

		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hash, err := hashPassword(body.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal encrypt"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"hashedPassword": hash,
		})
	})

	// ==================
	// ENDPOINT COMPARE
	// ==================
	r.POST("/compare", func(c *gin.Context) {
		var body struct {
			Password string `json:"password"`
			Hash     string `json:"hash"`
		}

		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		isValid := checkPassword(body.Password, body.Hash)

		c.JSON(http.StatusOK, gin.H{
			"success": isValid,
			"message": func() string {
				if isValid {
					return "Login berhasil"
				}
				return "Password salah"
			}(),
		})
	})

	r.Run(":8080")
}
