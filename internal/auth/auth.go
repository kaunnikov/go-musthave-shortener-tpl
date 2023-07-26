package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"kaunnikov/go-musthave-shortener-tpl/internal/errs"
	"kaunnikov/go-musthave-shortener-tpl/internal/logging"
	"net/http"
	"time"
)

const SecretKey = "UZo57ez$4e2V"
const CookieTokenName = "token"

type Claims struct {
	jwt.RegisteredClaims
	Token string
}

func GetUserToken(w http.ResponseWriter, r *http.Request) (string, error) {
	// получаем токен из куки
	tokenCookie, _ := r.Cookie(CookieTokenName)

	// Если токена нет - сформируем новый и запишем клиенту в куку
	if tokenCookie == nil {
		tokenCookie = generateCookie()
		http.SetCookie(w, tokenCookie)
	}

	// Достаём токен из куки и расшифровываем
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenCookie.Value, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(SecretKey), nil
		})

	if err != nil {
		logging.Errorf("Token don't decode: %s", err)
		return "", err
	}

	// Если кука не валидная - удаляем старую и пробуем снова
	if !token.Valid {
		deleteCoolie(w)
		//GetUserToken(w, r)
	}

	if claims.Token == "" {
		logging.Errorf("Token not found in cookie: %s", tokenCookie)
		return "", &errs.TokenNotFoundInCookie{
			Err: fmt.Errorf("token not found in cookie: %s", tokenCookie),
		}
	}

	return claims.Token, nil
}

func deleteCoolie(w http.ResponseWriter) {
	c := &http.Cookie{
		Name:    CookieTokenName,
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),

		HttpOnly: true,
	}

	http.SetCookie(w, c)
}

func generateCookie() *http.Cookie {
	token, err := generateJWTString()
	if err != nil {
		//
	}

	return &http.Cookie{
		Name:  "token",
		Value: token,
		Path:  "/",
	}
}

func generateJWTString() (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(31 * 24 * time.Hour)),
		},
		Token: uuid.NewString(),
	})

	return token.SignedString([]byte(SecretKey))
}
