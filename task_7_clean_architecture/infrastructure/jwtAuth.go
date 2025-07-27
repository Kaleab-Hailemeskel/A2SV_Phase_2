package infrastructure

import (
	"fmt"
	"log"
	"task_7_clean_architecture/models"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// ! Should Have it's own interface of Authenticatoin interface
type jWTAuth struct {
	JWTSECRET_TOKEN     string
	JWTSigningMethod    *jwt.SigningMethodHMAC
	TokenExpirationTime time.Duration
	HeaderName          string
}

func NewJWTAuth() models.IAuthentication {
	return &jWTAuth{
		JWTSECRET_TOKEN:     JWTSECRET,
		JWTSigningMethod:    jwt.SigningMethodHS256,
		TokenExpirationTime: time.Duration(2 * 24 * time.Hour), // ! Hard coded
		HeaderName:          HEADER,
	}
}

// ? This method shouldn't be accessed by outside folders
// sends empty string if generating failed
func (auth *jWTAuth) GenerateSecurityToken(JWTBody map[string]interface{}) (string, *time.Duration) {
	JWTBody["expiration_date"] = auth.TokenExpirationTime // set Expiration time
	token := jwt.NewWithClaims(auth.JWTSigningMethod, jwt.MapClaims(JWTBody))
	log.Println(auth.JWTSECRET_TOKEN)
	jwtToken, err := token.SignedString([]byte(auth.JWTSECRET_TOKEN))
	if err != nil {
		log.Printf("%s can not generate JWT\n", err)
		return "", nil
	}
	return jwtToken, &auth.TokenExpirationTime
}

// ? ABRHAM MENTIONED A WAY TO HACK JWT(a way to pass with nil) consider that after the first submission
func (auth *jWTAuth) GetUserEmailFromSecurityToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signging method %v", token.Header["alg"])
		}
		return []byte(auth.JWTSECRET_TOKEN), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) >= claims["expiration_date"].(float64) { //? REFRESH KEY should be considered HERE !!!
			return "", fmt.Errorf("expired JWT")
		}
		userEmail := claims["email"].(string)
		return userEmail, nil // valid return
	}
	return "", fmt.Errorf("jwtMapClaims raised an error")
}
