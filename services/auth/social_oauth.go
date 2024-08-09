package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"
	"gorm.io/gorm"

	"github.com/hngprojects/telex_be/internal/models"
	"github.com/hngprojects/telex_be/pkg/middleware"
	"github.com/hngprojects/telex_be/pkg/repository/storage"
	"github.com/hngprojects/telex_be/pkg/repository/storage/postgresql"
	"github.com/hngprojects/telex_be/services/actions"
	"github.com/hngprojects/telex_be/services/actions/names"
	"github.com/hngprojects/telex_be/utility"
)

func CreateGoogleUser(req models.GoogleRequestModel, db *gorm.DB) (gin.H, int, error) {

	var (
		userClaims   map[string]interface{}
		reqUser      models.CreateUserRequestModel
		sendWelcome  bool
		responseData gin.H
	)

	tokenString := req.Token

	resp, err := idtoken.Validate(context.Background(), tokenString, "")
	userClaims = resp.Claims
	if err != nil {
		return responseData, http.StatusBadRequest, fmt.Errorf("token not valid: " + err.Error())
	}

	var (
		email    = strings.ToLower(userClaims["email"].(string))
		username = strings.ToLower(userClaims["name"].(string))
		user     models.User
	)

	if email == "" || username == "" {
		return responseData, http.StatusNotFound, fmt.Errorf("token decode failed")
	}

	reqUser = models.CreateUserRequestModel{
		Email: email,
	}
	_, err = ValidateCreateUserRequest(reqUser, db)
	if err != nil {
		exists := postgresql.CheckExists(db, &user, "email = ?", email)
		if !exists {
			return responseData, http.StatusNotFound, fmt.Errorf("user not found")
		}
		user, err = user.GetUserWithProfile(db, user.ID)

		if err != nil {
			return responseData, http.StatusInternalServerError, fmt.Errorf("error fetching user " + err.Error())
		}

	} else {
		user = models.User{
			ID:         utility.GenerateUUID(),
			Name:       username,
			Email:      email,
			IsVerified: true,
			Profile: models.Profile{
				ID:        utility.GenerateUUID(),
				AvatarURL: userClaims["picture"].(string),
			},
		}
		err := user.CreateUser(db)
		sendWelcome = true
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
		"user": map[string]interface{}{
			"id":          user.ID,
			"email":       user.Email,
			"username":    user.Name,
			"fullname":    user.Name,
			"is_verified": user.IsVerified,
			"avatar_url":  user.Profile.AvatarURL,
			"expires_in":  strconv.Itoa(int(tokenData.ExpiresAt.Unix())),
			"created_at":  strconv.Itoa(int(user.CreatedAt.Unix())),
			"updated_at":  strconv.Itoa(int(user.UpdatedAt.Unix())),
		},
		"access_token": tokenData.AccessToken,
	}
	if sendWelcome {
		resetReq := models.SendWelcomeMail{
			Email: user.Email,
		}

		err = actions.AddNotificationToQueue(storage.DB.Redis, names.SendWelcomeMail, resetReq)
		if err != nil {
			return responseData, http.StatusInternalServerError, err
		}
	}

	return responseData, http.StatusCreated, nil
}
