package utils

import (
	 "errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secretKey="dummysecretkey"




func GenerateToken(email string, userId int64) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "email":  email,
        "userId": userId,
        "exp":    time.Now().Add(time.Hour * 2).Unix(), // expire time
    })
    return token.SignedString([]byte(secretKey))
}


func VerifyToken(token string) (int64, error) {
    parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return []byte(secretKey), nil
    })

    if err != nil {
        return 0, errors.New("could not parse token")
    }

    if !parsedToken.Valid {
        return 0, errors.New("invalid token")
    }

    claims, ok := parsedToken.Claims.(jwt.MapClaims)
    if !ok || !parsedToken.Valid {
        return 0, errors.New("invalid token claims")
    }

    userId, ok := claims["userId"].(float64)
    if !ok {
        return 0, errors.New("userId claim is missing or invalid")
    }

    return int64(userId), nil
}