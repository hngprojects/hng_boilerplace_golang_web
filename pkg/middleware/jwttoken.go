package middleware

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/config"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type TokenDetailDTO struct {
	AccessUuid  string `json:"access_uuid"`
	AccessToken string `json:"access_token"`
	ExpiresAt   time.Time
}

func CreateToken(user models.User) (*TokenDetailDTO, error) {

	var (
		tokenData = &TokenDetailDTO{}
		config    = config.GetConfig()
		err       error
	)

	tokenData.ExpiresAt = time.Now().AddDate(0, 0, config.Server.AccessTokenExpireDuration) // token valid for env set days
	tokenData.AccessUuid = user.ID
	tokenData.AccessUuid = utility.GenerateUUID()

	//create token
	userClaims := jwt.MapClaims{}

	// specify user claims
	userClaims["user_id"] = user.ID
	userClaims["access_uuid"] = tokenData.AccessUuid
	userClaims["exp"] = tokenData.ExpiresAt.Unix()
	userClaims["authorised"] = true

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)

	tokenData.AccessToken, err = token.SignedString([]byte(config.Server.Secret))
	if err != nil {
		return tokenData, err
	}

	return tokenData, nil
}

// verify token

func verifyToken(tokenString string) (*jwt.Token, error) {
	config := config.GetConfig()
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Server.Secret), nil
	})
	if err != nil {
		return token, fmt.Errorf("Unauthorized")
	}
	return token, nil
}

// check if token is valid
func TokenValid(bearerToken string) (*jwt.Token, error) {
	token, err := verifyToken(bearerToken)
	if err != nil {
		if token != nil {
			return token, err
		}
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("Unauthorized")
	}
	return token, nil
}

func GetUserClaims(c *gin.Context, db *gorm.DB, theValue string) (interface{}, error) {

	claims, exists := c.Get("userClaims")
	if !exists {
		return nil, errors.New("user claims not found")
	}

	userClaims := claims.(jwt.MapClaims)
	userValue, ok := userClaims[theValue]
	if !ok {
		return nil, errors.New("invalid value")
	}

	return userValue, nil

}
