package middleware

import (
    "context"
    "net/http"
    "strings"
    "github.com/golang-jwt/jwt/v5"
    "gorm.io/gorm"
    "noteslord/config"
    "noteslord/models"
)

type contextKey string

const userCtxKey = contextKey("userID")

func JWTAuthMiddleware(db *gorm.DB) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
                return
            }

            tokenString := strings.TrimPrefix(authHeader, "Bearer ")
            if tokenString == authHeader {
                http.Error(w, "Invalid token format", http.StatusUnauthorized)
                return
            }

            token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
                return []byte(config.JWTSecret), nil
            })
            if err != nil || !token.Valid {
                http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
                return
            }

            claims, ok := token.Claims.(jwt.MapClaims)
            if !ok {
                http.Error(w, "Invalid token claims", http.StatusUnauthorized)
                return
            }

            userIDFloat, ok := claims["user_id"].(float64)
            if !ok {
                http.Error(w, "Invalid token payload", http.StatusUnauthorized)
                return
            }

            userID := uint(userIDFloat)

            var user models.User
            if err := db.First(&user, userID).Error; err != nil {
                http.Error(w, "User not found", http.StatusUnauthorized)
                return
            }

            ctx := context.WithValue(r.Context(), userCtxKey, userID)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func GetUserIDFromContext(r *http.Request) (uint, bool) {
    userID, ok := r.Context().Value(userCtxKey).(uint)
    return userID, ok
}
