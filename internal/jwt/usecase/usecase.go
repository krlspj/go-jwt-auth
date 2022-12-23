package usecase

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JwtClaims struct {
	Username string `json:"username,omitemtpy"`
	UserId   string `json:"userId,omitempty"`
	jwt.StandardClaims
}

// type JwtUsecaser interface {
// JwtUsecase defines the jwt service
type JwtUsecase interface {
	GenerateToken(lifeTime int, userId, username string) (string, error)
	//ValidateToken(signedToken string) (jwtClaims, error)
	ValidateToken(signedToken string) (JwtClaims, error)
}

type defJwtUsecase struct {
	bSecret []byte
}

func NewJwtUsecase(hmacSampleSecret string) *defJwtUsecase {
	return &defJwtUsecase{
		bSecret: []byte(hmacSampleSecret),
	}
}

// generateToken returns the generated token and an error
func (s *defJwtUsecase) GenerateToken(lifeTime int, userId, username string) (string, error) {
	claims := &JwtClaims{
		Username: username,
		UserId:   userId,
		//Email:    email,
		//RoleId: role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Duration(lifeTime) * time.Minute).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.bSecret)

	//refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(secret_key))
	if err != nil {
		log.Panic(err)
		return "", err
	}

	return token, nil
}

func (s *defJwtUsecase) ValidateToken(signedToken string) (JwtClaims, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return s.bSecret, nil
		})

	if token.Valid {
		fmt.Println("You look nice today")
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		fmt.Println("That's not even a token")
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		// Token is either expired or not active yet
		fmt.Println("Timing is everything")
	} else {
		fmt.Println("Couldn't handle this token:", err)
	}

	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return *claims, nil

	} else {
		log.Println("[ERROR] [JWT] err:", err.Error(), "- Ok:", ok)
		return JwtClaims{}, err
	}
}

// --------- Token Public functions ----------
const tokenSecret string = "myTokenSecret"

// generateToken returns the generated token and an error
func GenerateToken(lifeTime int, userId, username string) (string, error) {
	claims := &JwtClaims{
		Username: username,
		UserId:   userId,
		//Email:    email,
		//RoleId: role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Duration(lifeTime) * time.Minute).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(tokenSecret)

	//refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(secret_key))
	if err != nil {
		log.Panic(err)
		return "", err
	}

	return token, nil
}

func ValidateToken(signedToken string) (JwtClaims, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return tokenSecret, nil
		})

	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return *claims, nil
	} else {
		log.Println("[ERROR] [JWT] [usecase] Error on claims", err, "Ok ->", ok)
		return JwtClaims{}, err
	}
}
