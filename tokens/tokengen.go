package tokens

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var SECRET_KEY = os.Getenv("SECRET_KEY")

type SignedDetails struct {
	Email      string
	First_Name string
	Last_Name  string
	Uid        string
	jwt.RegisteredClaims
}

// генерирует access и refresh токены
func TokenGenerator(email string, firstname string, lastname string) (signedToken string, signedRefreshToken string, err error) {
	if SECRET_KEY == "" {
		SECRET_KEY = "your-secret-key-change-this-in-production"
		log.Println("WARNING: Using default SECRET_KEY. Set SECRET_KEY environment variable in production!")
	}

	claims := &SignedDetails{
		Email:      email,
		First_Name: firstname,
		Last_Name:  lastname,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshClaims := &SignedDetails{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(168 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err = token.SignedString([]byte(SECRET_KEY))

	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	signedRefreshToken, err = refreshToken.SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)

		return "", "", err
	}

	return signedToken, signedRefreshToken, nil
}

// ValidateToken проверяет валидность токена
func ValidateToken(signedToken string) (claims *SignedDetails, err error) {
	if SECRET_KEY == "" {
		SECRET_KEY = "your-secret-key-change-this-in-production"
	}

	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},

		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*SignedDetails)

	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, fmt.Errorf("token is already expired")
	}

	return claims, nil
}

// обновляет токены пользователя в базе данных
func UpdateAllTokens(db *pgxpool.Pool, signedToken string, signedRefreshToken string, userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	updateQuery := `
		UPDATE users
		SET token = $1, refresh_token = $2, updated_at = $3
		WHERE user_id = $4
	`

	_, err := db.Exec(ctx, updateQuery, signedToken, signedRefreshToken, time.Now().UTC(), userId)

	if err != nil {
		log.Printf("Error updating tokens for user %s: %v", userId, err)
		return err
	}

	return nil
}
