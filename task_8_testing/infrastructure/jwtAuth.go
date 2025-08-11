package infrastructure

import (
	"fmt"
	"log"
	"task_8_testing/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ! Should Have it's own interface of Authenticatoin interface
type jWTAuth struct {
	JWTSECRET_TOKEN     string
	JWTSigningMethod    *jwt.SigningMethodHMAC
	TokenExpirationTime time.Duration
}

func NewJWTAuth() models.IAuthentication {
	return &jWTAuth{
		JWTSECRET_TOKEN:     JWTSECRET,
		JWTSigningMethod:    jwt.SigningMethodHS256,
		TokenExpirationTime: time.Duration(24 * time.Hour), // ! Hard coded
	}
}

// ? This method shouldn't be accessed by outside folders
// sends empty string if generating failed
func (auth *jWTAuth) GenerateSecurityToken(JWTBody map[string]interface{}) (string, time.Duration) {
	JWTBody["expiration_date"] = auth.TokenExpirationTime // set Expiration time
	token := jwt.NewWithClaims(auth.JWTSigningMethod, jwt.MapClaims(JWTBody))
	log.Println(auth.JWTSECRET_TOKEN)
	jwtToken, err := token.SignedString([]byte(auth.JWTSECRET_TOKEN))
	if err != nil {
		log.Printf("%s can not generate JWT\n", err)
		return "", time.Duration(0)
	}
	return jwtToken, auth.TokenExpirationTime
}
func (auth *jWTAuth) TokenExpired(token *jwt.Token) (bool, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) >= claims["expiration_date"].(float64) {
			return true, fmt.Errorf("expired JWT")
		}
		return false, nil // valid return
	}
	return false, fmt.Errorf("token don't have expiration date")
}

func (auth *jWTAuth) ParseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signging method %v", token.Header["alg"])
		}
		return []byte(auth.JWTSECRET_TOKEN), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// ? ABRHAM MENTIONED A WAY TO HACK JWT(a way to pass with nil) consider that after the first submission
func (auth *jWTAuth) GetUserEmailFromSecurityToken(token *jwt.Token) (string, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userEmail := claims["email"].(string)
		return userEmail, nil // valid return
	}
	return "", fmt.Errorf("jwtMapClaims raised an error")
}
func (auth *jWTAuth) GetUserID(token *jwt.Token) (*primitive.ObjectID, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userID := claims["id"].(primitive.ObjectID)
		return &userID, nil // valid return
	}
	return nil, fmt.Errorf("jwtMapClaims raised an error")
}
