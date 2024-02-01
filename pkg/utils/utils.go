package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var secretKey = []byte("No_Ra1nB0wT4bl3_c4n_H3lp!")

func JWTGenerateToken(claims map[string]interface{}) (string, error) {
	// Append the expiration claim (4 minutes)
	claims["exp"] = time.Now().Add(time.Minute * 90).Unix()

	// Create a new token object, specifying signing method and the claims you want to include
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))

	// Sign the token with a secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
func JWTDecodeToken(tokenString string) (map[string]interface{}, bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method and return the secret key
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, false, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, false, fmt.Errorf("invalid token")
	}

	// Check if the token is expired
	expired := claims.VerifyExpiresAt(time.Now().Unix(), false)

	// Remove the "exp" field from the claims
	delete(claims, "exp")

	return claims, expired, nil
}
func GenerateRefreshToken(userID string) (string, error) {
	// Set the expiration time for the refresh token to one year
	expirationTime := time.Now().Add(365 * 24 * time.Hour)

	// Create a new refresh token object with user-specific claims
	refreshTokenClaims := jwt.MapClaims{
		"userid": userID,                // Subject (user ID)
		"exp":    expirationTime.Unix(), // Expiration time
		"iss":    "API Services Center", // Issuer
		// Add other user-specific claims as needed
	}

	// Create a new refresh token, specifying the signing method and claims
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	// Sign the refresh token with a secret key
	refreshTokenString, err := refreshToken.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return refreshTokenString, nil
}
func isValidRefreshToken(refreshToken string) (bool, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return false, err
	}

	// Check if the token is not expired and has a valid signature
	return token.Valid, nil
}
