package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

type CustomClaim struct {
	User string `json:"username,omitempty"`
	jwt.StandardClaims
}

//GenerateJWT with user payload for auth validation
func GenerateJWT(username string) string{
	signingKey := []byte(os.Getenv("sign_key"))

	claims :=  CustomClaim{
		User: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
			IssuedAt: time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(signingKey)
	return tokenString
}

func GenerateRefreshToken() string{
	signingKey := []byte(os.Getenv("refresh_key"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().AddDate(0,1,0).Unix(),
	})

	tokenString, _ := token.SignedString(signingKey)
	return tokenString
}

func ValidateToken(refToken string, reqToken string) (success bool, username string, token string) {
	signRefKey := []byte(os.Getenv("refresh_key"))
	signKey := []byte(os.Getenv("sign_key"))
	//check refresh token first
	if refToken == ""{
		return false,"",""
	}
	tokenRef, err := jwt.Parse(refToken,func(token *jwt.Token) (interface{}, error){
		return signRefKey, nil
	})
	if !tokenRef.Valid || err != nil{
		return false, "", ""
	}
	claims := &CustomClaim{}
	jwt.ParseWithClaims(reqToken, claims, func(token *jwt.Token)(interface{}, error){
		return signKey, nil
	})
	//refresh user token
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 5*time.Minute{
		newToken := GenerateJWT(claims.User)
		return true, claims.User, newToken
	}
	return true, claims.User, reqToken
}