package helpers

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func SetupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, x-api-key, token")
}

func HashPassword(password string) (string, error) { //Hash a users password before storing in database
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool { //check a users password with hashed password, usually for login purpose
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
