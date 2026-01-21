package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

//
// =======================
// FUNGSI BCRYPT
// =======================
//

// Hash password (dipakai saat membuat user / register)
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	return string(bytes), err
}

// Mengecek password saat login
func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
	return err == nil
}

//
// =======================
// STRUKTUR DATA USER
// =======================
//

// Struktur user (simulasi tabel user)
type User struct {
	Username string
	Hash     string // password yang sudah di-hash
}

// Data user (simulasi database)
var users = []User{}

func main() {
	// =======================
	// DATA CONTOH USER
	// =======================
	// Username: admin
	// Password: password
	hash, _ := hashPassword("password")
	users = append(users, User{
		Username: "admin",
		Hash:     hash,
	})

	// Inisialisasi Gin
	r := gin.Default()

	// CORS agar bisa diakses dari React
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"POST"},
		AllowHeaders: []string{"Content-Type"},
	}))

	// =======================
	// ENDPOINT LOGIN
	// =======================
	r.POST("/", func(c *gin.Context) {
		// Struktur body request
		var body struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		// Validasi JSON
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Format request tidak valid",
			})
			return
		}

		// Validasi input kosong
		if body.Username == "" || body.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Username dan password wajib diisi",
			})
			return
		}

		// Cari user berdasarkan username
		var found *User
		for i := range users {
			if users[i].Username == body.Username {
				found = &users[i]
				break
			}
		}

		// Jika username tidak ditemukan
		if found == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Username tidak ditemukan",
			})
			return
		}

		// Cek password dengan bcrypt
		if !checkPassword(body.Password, found.Hash) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Password salah",
			})
			return
		}

		// Login berhasil
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Login berhasil",
		})
	})

	// Jalankan server
	r.Run(":8080")
}
