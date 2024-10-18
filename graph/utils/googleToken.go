package utils

import (
	crand "crypto/rand"
	"encoding/base64"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var GoogleClientIDWeb string
var GoogleClientIDIos string
var GoogleClientIDAndroid string

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	GoogleClientIDWeb = os.Getenv("GOOGLE_CLIENT_ID_WEB")
	GoogleClientIDIos = os.Getenv("GOOGLE_CLIENT_ID_IOS")
	GoogleClientIDAndroid = os.Getenv("GOOGLE_CLIENT_ID_ANDROID")
}
func GenerateRandomString(n int) (string, error) {
	b := make([]byte, n)
	_, err := crand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
