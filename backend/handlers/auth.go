package handlers

import (
    "encoding/json"
    "net/http"
    "regexp"
    "strings"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"

    "your_project/models"
)

var jwtKey = []byte("your_secret_key")

type RegisterInput struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

type LoginInput struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

func isValidEmail(email string) bool { //this is a helper function to validate email format using regex
	//regex pattern for validating email addresses
    regex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
    return regex.MatchString(strings.ToLower(email))
}


func Register(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var input RegisterInput
        if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
            http.Error(w, "Invalid input format", http.StatusBadRequest)
            return
        }

        input.Username = strings.TrimSpace(input.Username)
        input.Email = strings.TrimSpace(input.Email)
        if input.Username == "" || input.Email == "" || input.Password == "" {
            http.Error(w, "All fields are required", http.StatusBadRequest)
            return
        }

        if !isValidEmail(input.Email) {
            http.Error(w, "Invalid email format", http.StatusBadRequest)
            return
        }

        if len(input.Password) < 6 {
            http.Error(w, "Password must be at least 6 characters", http.StatusBadRequest)
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

        // Basic validation
        if input.Email == "" || input.Password == "" {
            http.Error(w, "Email and password are required", http.StatusBadRequest)
            return
        }

        if !isValidEmail(input.Email) {
            http.Error(w, "Invalid email format", http.StatusBadRequest)
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

        token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
            "user_id": user.ID,
            "exp":     time.Now().Add(24 * time.Hour).Unix(),
        })

        tokenString, err := token.SignedString(jwtKey)
        if err != nil {
            http.Error(w, "Could not generate token", http.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
    }
}
