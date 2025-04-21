package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte(os.Getenv("TODO_PASSWORD"))

type claims struct {
	jwt.RegisteredClaims
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	pass := os.Getenv("TODO_PASSWORD")
	if pass == "" {
		writeJSON(w, map[string]string{"error": "authentication disabled"})
		return
	}
	if req.Password != pass {
		writeJSON(w, map[string]string{"error": "Неверный пароль"})
		return
	}
	now := time.Now().UTC()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(8 * time.Hour)),
			Subject:   "todo-user",
		},
	})
	signed, err := token.SignedString(jwtKey)
	if err != nil {
		writeJSON(w, map[string]string{"error": "token generation error"})
		return
	}
	writeJSON(w, map[string]string{"token": signed})
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if pass != "" {
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}
			tok, err := jwt.ParseWithClaims(cookie.Value, &claims{}, func(t *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})
			if err != nil || !tok.Valid {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	}
}
