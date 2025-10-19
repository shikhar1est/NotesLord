package handlers

import (
    "encoding/json"
    "net/http"
    "strings"
    "time"
    "github.com/go-playground/validator/v10"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
    "noteslord/config"
    "noteslord/models"
)

var validate = validator.New()

type RegisterInput struct {
    Username string `json:"username" validate:"required,min=3,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

type LoginInput struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

func Register(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var input RegisterInput
        if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
            http.Error(w, "Invalid input format", http.StatusBadRequest)
            return
        }

        // Trim extra spaces
        input.Username = strings.TrimSpace(input.Username)
        input.Email = strings.TrimSpace(input.Email)

        // Validate struct
        if err := validate.Struct(input); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        // Hash password
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
        if err != nil {
            http.Error(w, "Error hashing password", http.StatusInternalServerError)
            return
        }

        user := models.User{
            Username: input.Username,
            Email:    input.Email,
            Password: string(hashedPassword),
        }

        if err := db.Create(&user).Error; err != nil {
            http.Error(w, "User already exists or invalid data", http.StatusBadRequest)
            return
        }

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
    }
}

// 🟨 Login Handler
func Login(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var input LoginInput
        if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
            http.Error(w, "Invalid input format", http.StatusBadRequest)
            return
        }

        input.Email = strings.TrimSpace(input.Email)

        if err := validate.Struct(input); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        var user models.User
        if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
            http.Error(w, "Invalid credentials", http.StatusUnauthorized)
            return
        }

        if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
            http.Error(w, "Invalid credentials", http.StatusUnauthorized)
            return
        }

        // Generate JWT Token
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
            "user_id": user.ID,
            "exp":     time.Now().Add(24 * time.Hour).Unix(),
        })

        tokenString, err := token.SignedString([]byte(config.JWTSecret))
        if err != nil {
            http.Error(w, "Could not generate token", http.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
    }
}
