package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	password string
	jwtKey   []byte
)

func init() {
	password = os.Getenv("TODO_PASSWORD")
	jwtKey = []byte(password)
}

type claims struct {
	jwt.RegisteredClaims
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if password != "" {
			cookie, err := r.Cookie("token")
			if err != nil {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "authentication required"})
				return
			}
			tok, err := jwt.ParseWithClaims(cookie.Value, &claims{}, func(t *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})
			if err != nil || !tok.Valid {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "authentication required"})
				return
			}
		}
		next(w, r)
	}
}
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	if password == "" {
		writeJSON(w, http.StatusOK, map[string]string{"message": "authentication disabled"})
		return
	}
	if req.Password != password {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid password"})
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
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "token generation error"})
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    signed,
		Expires:  now.Add(8 * time.Hour),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})
	writeJSON(w, http.StatusOK, map[string]string{"token": signed})
}
