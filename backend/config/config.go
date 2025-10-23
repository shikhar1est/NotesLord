package config
//config file
import (
    "log"
    "os"
    "github.com/joho/godotenv"
)

var JWTSecret string

func LoadEnv() {
    err := godotenv.Load()
    if err != nil {
        log.Println("Warning: .env file not found")
    }

    JWTSecret = os.Getenv("JWT_SECRET")
    if JWTSecret == "" {
        log.Fatal("JWT_SECRET not found in environment variables")
    }
}
