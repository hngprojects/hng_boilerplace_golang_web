package token

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/config"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
)

func GetConnToken(userId string, db *gorm.DB) (gin.H, int, error) {

	userClaims := jwt.MapClaims{}

	userClaims["sub"] = userId
	userClaims["exp"] = time.Now().Unix() + int64(120)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)

	connToken, err := token.SignedString([]byte(config.Config.Centrifuge.Secret))
	if err != nil {
		return gin.H{}, http.StatusInternalServerError, err
	}

	res := gin.H{
		"token": connToken,
	}

	return res, http.StatusOK, nil

}

func GetSubToken(userId string, req models.ChannelSubTokenReq, db *gorm.DB) (gin.H, int, error) {

	var (
		channelName = req.Channel
	)

	userClaims := jwt.MapClaims{}

	userClaims["sub"] = userId
	userClaims["channel"] = channelName
	userClaims["exp"] = time.Now().Unix() + int64(300)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)

	subToken, err := token.SignedString([]byte(config.Config.Centrifuge.Secret))
	if err != nil {
		return gin.H{}, http.StatusInternalServerError, err
	}

	res := gin.H{
		"token": subToken,
	}

	return res, http.StatusOK, nil
}
