package handlers

import (
    "encoding/json"
    "net/http"

    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
    "noteslord/models"
)

type RegisterInput struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

func Register(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var input RegisterInput
        if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
            http.Error(w, "Invalid input", http.StatusBadRequest)
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
            http.Error(w, "User already exists", http.StatusBadRequest)
            return
        }

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
    }
}
