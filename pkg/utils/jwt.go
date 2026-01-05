package utils

import (
	"errors"
	"html"
	"net/http"
	"strings"
	"time"

	"github.com/adohong4/driving-license/config"
	"github.com/adohong4/driving-license/internal/models"
	"github.com/dgrijalva/jwt-go"
)

// JWT Claims struct
type Claims struct {
	IdentityNO string `json:"identity_no"`
	ID         string `json:"id"`
	jwt.StandardClaims
}

type ClaimsUserAddress struct {
	UserAddress string `json:"user_address"`
	IdentityNO  string `json:"identity_no"`
	ID          string `json:"id"`
	jwt.StandardClaims
}

type ClaimsGovAgencyAddress struct {
	UserAddress string `json:"user_address"`
	jwt.StandardClaims
}

// Generate new JWT Token
func GenerateJWTToken(user *models.User, config *config.Config) (string, error) {
	// Register the JWT claims, which includes the username and expiry time
	claims := &Claims{
		IdentityNO: user.IdentityNo,
		ID:         user.Id.String(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 60).Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// register the JWT string
	tokenString, err := token.SignedString([]byte(config.Server.JwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Generate new JWT Token from user_address
func GenerateJWTTokenFromUserAddress(user *models.User, config *config.Config) (string, error) {
	// Register the JWT claims, which includes the username and expiry time
	claims := &ClaimsUserAddress{
		UserAddress: *user.UserAddress,
		IdentityNO:  user.IdentityNo,
		ID:          user.Id.String(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 60).Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// register the JWT string
	tokenString, err := token.SignedString([]byte(config.Server.JwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateJWTTokenFromAgencyAddress(agency *models.GovAgency, config *config.Config) (string, error) {
	claims := &ClaimsGovAgencyAddress{
		UserAddress: agency.UserAddress,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 60).Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// register the JWT string
	tokenString, err := token.SignedString([]byte(config.Server.JwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Extract JWT From Request
func ExtractJWTFromRequest(r *http.Request) (map[string]interface{}, error) {
	//Get the JWT string
	tokenString := ExtractBearerToken(r)

	// Initialize a new instance of `Claims` (here using Claims map)
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (jwtKey interface{}, err error) {
		return jwtKey, err
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, errors.New("Invalid token signature")
		}
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("Invalid Token")
	}

	return claims, nil
}

// Extract bearer token from request Authorization header
func ExtractBearerToken(r *http.Request) string {
	headerAuthorization := r.Header.Get("Authorization")
	bearerToken := strings.Split(headerAuthorization, " ")
	return html.EscapeString(bearerToken[1])
}
