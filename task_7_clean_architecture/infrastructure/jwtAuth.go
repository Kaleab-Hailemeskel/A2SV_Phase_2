package infrastructure

import (
	"fmt"
	"log"
	"net/http"
	"task_7_clean_architecture/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// ! Should Have it's own interface of Authenticatoin interface
type JWTAuth struct {
	JWTSECRET_TOKEN     string
	JWTSigningMethod    *jwt.SigningMethodHMAC
	TokenExpirationTime time.Duration
	HeaderName          string
}

func NewJWTAuth() models.IAuthentication {
	return &JWTAuth{
		JWTSECRET_TOKEN:     JWTSECRET,
		JWTSigningMethod:    jwt.SigningMethodHS256,
		TokenExpirationTime: time.Duration(2 * 24 * time.Hour), // ! Hard coded
		HeaderName:          HEADER,
	}
}

// ? This method shouldn't be accessed by outside folders
// sends empty string if generating failed
func (auth *JWTAuth) generateToken(JWTBody map[string]interface{}) string {
	JWTBody["expiration_date"] = auth.TokenExpirationTime // set Expiration time
	token := jwt.NewWithClaims(auth.JWTSigningMethod, jwt.MapClaims(JWTBody))
	log.Println(auth.JWTSECRET_TOKEN)
	jwtToken, err := token.SignedString([]byte(auth.JWTSECRET_TOKEN))
	if err != nil {
		log.Printf("%s can not generate JWT\n", err)
		return ""
	}
	return jwtToken
}
func (auth *JWTAuth) SendSecurityTokenToClinet(ctx *gin.Context, JWTBody map[string]interface{}) error {
	jwtToken := auth.generateToken(JWTBody)
	if jwtToken == "" {
		return fmt.Errorf("empty JWT")
	}
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(auth.HeaderName, jwtToken, int(time.Now().Add(auth.TokenExpirationTime).Unix()), "", "", false, true) // the int(time.Now().Add(auth.TokenExpirationTime).Unix()) part could be a field of the jwtAuth structure

	//  auth.JWTctx.JSON(200, gin.H{"message": "Cookies were sent"}) //? For Debugging purposes only

	return nil
}

// ? ABRHAM MENTIONED A WAY TO HACK JWT(a way to pass with nil) consider that after the first submission
func (auth *JWTAuth) GetSecurityTokenFromClinet(ctx *gin.Context) (string, error) {
	tokenString, err := ctx.Cookie(auth.HeaderName)
	if err != nil {
		return "", err
	}
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
		return userEmail, nil
	}
	return "", fmt.Errorf("jwtMapClaims raised an error")
}
