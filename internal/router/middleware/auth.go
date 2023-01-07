package middleware

import (
	"encoding/hex"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

func VerifyJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Header["Authorization"] != nil {
			token, err := jwt.Parse(request.Header["Authorization"][0], func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					writer.WriteHeader(http.StatusUnauthorized)
					_, err := writer.Write([]byte("You're Unauthorized"))
					if err != nil {
						return nil, err
					}
				}
				secretKey, err := hex.DecodeString(os.Getenv("SECRET_KEY"))
				if err != nil {
					return nil, err
				}
				return secretKey, nil

			})
			if err != nil {
				http.Error(writer, err.Error(), http.StatusUnauthorized)
				return
			}
			if token.Valid {
				next.ServeHTTP(writer, request)
			} else {
				writer.WriteHeader(http.StatusUnauthorized)
				_, err := writer.Write([]byte("You're Unauthorized due to invalid token"))
				if err != nil {
					return
				}
			}
		} else {
			writer.WriteHeader(http.StatusUnauthorized)
			_, err := writer.Write([]byte("You're Unauthorized due to No token in the header"))
			if err != nil {
				return
			}
		}
	})
}
