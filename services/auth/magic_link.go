package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/hngprojects/telex_be/internal/config"
	"github.com/hngprojects/telex_be/internal/models"
	"github.com/hngprojects/telex_be/pkg/middleware"
	"github.com/hngprojects/telex_be/pkg/repository/storage"
	"github.com/hngprojects/telex_be/pkg/repository/storage/postgresql"
	"github.com/hngprojects/telex_be/services/actions"
	"github.com/hngprojects/telex_be/services/actions/names"
	"github.com/hngprojects/telex_be/utility"
)

func MagicLinkRequest(userEmail, url string, db *gorm.DB) (string, int, error) {

	var (
		user      = models.User{}
		magicLink = models.MagicLink{}
		config    = config.GetConfig()
	)

	magicExist, err := magicLink.GetMagicLinkByEmail(db, userEmail)
	if err != nil {
		return "error", http.StatusUnauthorized, err
	}

	if magicExist != nil {
		if err := magicExist.DeleteMagicLink(db); err != nil {
			return "error", http.StatusInternalServerError, err
		}
	}

	exists := postgresql.CheckExists(db, &user, "email = ?", userEmail)
	if !exists {
		return "error", http.StatusNotFound, fmt.Errorf("user not found")
	}

	resetToken, err := utility.GenerateOTP(6)

	if err != nil {
		return "error", http.StatusInternalServerError, err
	}

	magic := models.MagicLink{
		ID:        utility.GenerateUUID(),
		Email:     strings.ToLower(userEmail),
		Token:     strconv.Itoa(resetToken),
		ExpiresAt: time.Now().Add(time.Duration(config.App.MagicLinkDuration) * time.Minute),
	}

	err = magic.CreateMagicLink(db)
	if err != nil {
		return "error", http.StatusInternalServerError, err
	}

	magic_link := fmt.Sprintf("%v/login/magic-link?token=%v", url, resetToken)

	resetReq := models.SendMagicLink{
		Email:     userEmail,
		MagicLink: magic_link,
	}

	err = actions.AddNotificationToQueue(storage.DB.Redis, names.SendMagicLink, resetReq)
	if err != nil {
		return "error", http.StatusInternalServerError, err
	}

	return "success", http.StatusOK, nil
}

func VerifyMagicLinkToken(req models.VerifyMagicLinkRequest, db *gorm.DB) (gin.H, int, error) {

	var (
		user         = models.User{}
		responseData gin.H
		magicLink    = models.MagicLink{}
	)

	magicExist, err := magicLink.GetMagicLinkByToken(db, req.Token)
	if err != nil {
		return responseData, http.StatusUnauthorized, errors.New("invalid or expired token")
	}

	exists := postgresql.CheckExists(db, &user, "email = ?", magicExist.Email)
	if !exists {
		return responseData, http.StatusBadRequest, errors.New("invalid credentials")
	}

	userData, err := user.GetUserByEmail(db, magicExist.Email)
	if err != nil {
		return responseData, http.StatusInternalServerError, errors.New("unable to fetch user")
	}

	tokenData, err := middleware.CreateToken(user)
	if err != nil {
		return responseData, http.StatusInternalServerError, errors.New("error saving token")
	}

	tokens := map[string]string{
		"access_token": tokenData.AccessToken,
		"exp":          strconv.Itoa(int(tokenData.ExpiresAt.Unix())),
	}

	access_token := models.AccessToken{ID: tokenData.AccessUuid, OwnerID: user.ID}

	err = access_token.CreateAccessToken(db, tokens)

	if err != nil {
		return responseData, http.StatusInternalServerError, errors.New("error saving token")
	}

	if err := magicExist.DeleteMagicLink(db); err != nil {
		return responseData, http.StatusInternalServerError, err
	}

	responseData = gin.H{

		"user": map[string]string{
			"id":         userData.ID,
			"email":      userData.Email,
			"username":   userData.Name,
			"first_name": userData.Profile.FirstName,
			"last_name":  userData.Profile.LastName,
			"fullname":   userData.Profile.FirstName + " " + userData.Profile.LastName,
			"phone":      userData.Profile.Phone,
			"expires_in": strconv.Itoa(int(tokenData.ExpiresAt.Unix())),
			"created_at": strconv.Itoa(int(userData.CreatedAt.Unix())),
			"updated_at": strconv.Itoa(int(userData.UpdatedAt.Unix())),
		},
		"access_token": tokenData.AccessToken,
	}

	return responseData, http.StatusOK, nil
}
