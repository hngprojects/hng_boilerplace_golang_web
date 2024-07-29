package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func CreateGoogleUser(req models.GoogleRequestModel, db *gorm.DB) (gin.H, int, error) {

	var userClaims models.GoogleClaims
	var reqUser models.CreateUserRequestModel

	tokenString := req.Token

	// Parse the token

	_, err := jwt.ParseWithClaims(tokenString, &userClaims, func(token *jwt.Token) (interface{}, error) {
		// Provide the key for HMAC verification (replace with your actual key)
		return []byte(""), nil
	})

	var (
		email        = strings.ToLower(userClaims.Email)
		username     = strings.ToLower(userClaims.Name)
		responseData gin.H
		user         models.User
	)

	reqUser = models.CreateUserRequestModel{
		Email: email,
	}

	// check if user already exists
	_, err = ValidateCreateUserRequest(reqUser, db)
	if err != nil && errors.Is(err, errors.New("user already exists with the given email")) {

		exists := postgresql.CheckExists(db, &user, "email = ?", email)
		if !exists {
			return responseData, http.StatusNotFound, fmt.Errorf("user not found")
		}

	} else {

		user = models.User{
			ID:    utility.GenerateUUID(),
			Name:  username,
			Email: email,
			Role:  int(models.RoleIdentity.User),
			Profile: models.Profile{
				ID:        utility.GenerateUUID(),
				AvatarURL: userClaims.Picture,
			},
		}
		err := user.CreateUser(db)
		if err != nil {
			return responseData, http.StatusInternalServerError, err
		}
	}

	tokenData, err := middleware.CreateToken(user)
	if err != nil {
		return responseData, http.StatusInternalServerError, fmt.Errorf("error saving token: " + err.Error())
	}

	tokens := map[string]string{
		"access_token": tokenData.AccessToken,
		"exp":          strconv.Itoa(int(tokenData.ExpiresAt.Unix())),
	}

	access_token := models.AccessToken{ID: tokenData.AccessUuid, OwnerID: user.ID}

	err = access_token.CreateAccessToken(db, tokens)

	if err != nil {
		return responseData, http.StatusInternalServerError, fmt.Errorf("error saving token: " + err.Error())
	}

	responseData = gin.H{
		"email": user.Email,
		"name":  user.Name,
		"id":    user.ID,
		"role":  models.UserRoleName,
		"profile": map[string]string{
			"first_name": user.Name,
			"avatar_url": user.Profile.AvatarURL,
			"expires_in": strconv.Itoa(int(tokenData.ExpiresAt.Unix())),
		},
		"access_token": tokenData.AccessToken,
	}

	return responseData, http.StatusCreated, nil
}
