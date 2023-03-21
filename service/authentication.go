package service

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// JwtWrapper wraps the signing key and the issuer
type JwtWrapperUser struct {
	SecretKey       string
	Issuer          string
	ExpirationHours int64
}

// JwtClaim adds email as a claim to the token
type JwtClaimUser struct {
	Email string
	jwt.StandardClaims
}

type JwtWrapperAdmin struct {
	SecretKey       string
	Issuer          string
	ExpirationHours int64
}

type JwtClaimAdmin struct {
	Email string
	jwt.StandardClaims
}

// Generate Token generates a jwt token
func (j *JwtWrapperUser) GenerateTokenUser(email string) (signedToken string, err error) {
	claims := &JwtClaimUser{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(j.ExpirationHours)).Unix(),
			Issuer:    j.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err = token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return
	}

	return
}

// Validate Token validates the jwt token // ยืนยันว่าหมดเวลารึยัง
func (j *JwtWrapperUser) ValidateTokenUser(signedToken string) (claims *JwtClaimUser, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaimUser{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.SecretKey), nil
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*JwtClaimUser)
	if !ok {
		err = errors.New("Couldn't parse claims")
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("JWT is expired")
		return
	}

	return

}

// Generate Token generates a jwt token
func (j *JwtWrapperAdmin) GenerateTokenAdmin(email string) (signedToken string, err error) {
	claims := &JwtClaimAdmin{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(j.ExpirationHours)).Unix(),
			Issuer:    j.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err = token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return
	}

	return
}

// Validate Token validates the jwt token // ยืนยันว่าหมดเวลารึยัง
func (j *JwtWrapperAdmin) ValidateTokenAdmin(signedToken string) (claims *JwtClaimAdmin, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaimAdmin{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.SecretKey), nil
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*JwtClaimAdmin)
	if !ok {
		err = errors.New("Couldn't parse claims")
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("JWT is expired")
		return
	}

	return

}
