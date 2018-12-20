package main

import (
	"log"
	"time"

	pb "github.com/aaronflower/shippy-service-user/proto/auth"
	jwt "github.com/dgrijalva/jwt-go"
)

var (
	key = []byte("Gakkiyui")
)

// CustomClaims is our custom metadata, which will be hashed
// and set as the second segment in or JWT
type CustomClaims struct {
	User *pb.User
	jwt.StandardClaims
}

// Authable defines an authorization interface.
type Authable interface {
	Decode(token string) (*CustomClaims, error)
	Encode(user *pb.User) (string, error)
}

// TokenService implements the Authable interface.
type TokenService struct {
	repo Repository
}

// Decode a token string into a token object.
func (s *TokenService) Decode(token string) (*CustomClaims, error) {
	// Parse the token
	tokenType, err := jwt.ParseWithClaims(
		token,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return key, nil
		},
	)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Validate the token and return the custom claims
	if claims, ok := tokenType.Claims.(*CustomClaims); ok && tokenType.Valid {
		return claims, nil
	}
	return nil, err
}

// Encode a claim into a JWT
func (s *TokenService) Encode(user *pb.User) (string, error) {
	expireToken := time.Now().Add(time.Hour * 72).Unix()
	// create teh claims
	claims := CustomClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "go.micro.srv.user",
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token and return
	str, err := token.SignedString(key)
	if err != nil {
		log.Fatal(err)
	}
	return str, nil
}
