package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

var (
	ApplicationName         string
	ServerPort              string
	LoginExpirationDuration time.Duration
	JwtSigningMethod        jwt.SigningMethod
	JwtSignatureKey         []byte
)

func InitConfig() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Read Application Name (default: "Payment API App")
	ApplicationName = getEnv("APPLICATION_NAME", "Payment API App")

	// Read Server Port (default: 8081)
	ServerPort = getEnv("SERVER_PORT", "8081")

	// Read Login Expiration Duration (default: 5 minutes)
	expirationStr := getEnv("LOGIN_EXPIRATION_DURATION", "5")
	expiration, err := strconv.Atoi(expirationStr)
	if err != nil {
		log.Fatalf("Failed to parse LOGIN_EXPIRATION_DURATION: %v", err)
	}
	LoginExpirationDuration = time.Duration(expiration) * time.Minute

	// Set JwtSigningMethod (default: HS256)
	JwtSigningMethod = jwt.SigningMethodHS256

	// Read Jwt Signature Key (default: "secret")
	JwtSignatureKey = []byte(getEnv("JWT_SIGNATURE_KEY", "secret"))
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
