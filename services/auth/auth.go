package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func ValidateCreateUserRequest(req models.CreateUserRequestModel, db *gorm.DB) (models.CreateUserRequestModel, error) {

	user := models.User{}
	profile := models.Profile{}

	if req.Email != "" {
		req.Email = strings.ToLower(req.Email)
		formattedMail, checkBool := utility.EmailValid(req.Email)
		if !checkBool {
			return req, fmt.Errorf("email address is invalid")
		}
		req.Email = formattedMail
		exists := postgresql.CheckExists(db, &user, "email = ?", req.Email)
		if exists {
			return req, errors.New("user already exists with the given email")
		}
	}

	if req.PhoneNumber != "" {
		req.PhoneNumber = strings.ToLower(req.PhoneNumber)
		phone, _ := utility.PhoneValid(req.PhoneNumber)
		req.PhoneNumber = phone
		exists := postgresql.CheckExists(db, &profile, "phone = ?", req.PhoneNumber)
		if exists {
			return req, errors.New("user already exists with the given phone")
		}

	}

	return req, nil
}

func GetUser(userIDStr string, db *gorm.DB) (models.User, error) {
	var userResp models.User

	userResp, err := userResp.GetUserByID(db, userIDStr)
	if err != nil {
		return userResp, err
	}

	return userResp, nil
}

func CreateUser(req models.CreateUserRequestModel, db *gorm.DB) (gin.H, int, error) {

	var (
		email        = strings.ToLower(req.Email)
		firstName    = strings.Title(strings.ToLower(req.FirstName))
		lastName     = strings.Title(strings.ToLower(req.LastName))
		username     = strings.ToLower(req.UserName)
		phoneNumber  = req.PhoneNumber
		password     = req.Password
		responseData gin.H
	)

	password, err := utility.HashPassword(req.Password)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	user := models.User{
		ID:       utility.GenerateUUID(),
		Name:     username,
		Email:    email,
		Password: password,
		Profile: models.Profile{
			ID:        utility.GenerateUUID(),
			FirstName: firstName,
			LastName:  lastName,
			Phone:     phoneNumber,
		},
	}

	err = user.CreateUser(db)
	if err != nil {
		return nil, http.StatusInternalServerError, err
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
		"user": map[string]string{
			"id":         user.ID,
			"email":      user.Email,
			"username":   user.Name,
			"first_name": user.Profile.FirstName,
			"last_name":  user.Profile.LastName,
			"fullname":   user.Profile.FirstName + " " + user.Profile.LastName,
			"phone":      user.Profile.Phone,
			"expires_in": strconv.Itoa(int(tokenData.ExpiresAt.Unix())),
			"created_at": strconv.Itoa(int(user.CreatedAt.Unix())),
			"updated_at": strconv.Itoa(int(user.UpdatedAt.Unix())),
		},
		"access_token": tokenData.AccessToken,
	}

	return responseData, http.StatusCreated, nil
}

func LoginUser(req models.LoginRequestModel, db *gorm.DB) (gin.H, int, error) {

	var (
		user         = models.User{}
		responseData gin.H
	)

	// Check if the user email exists
	exists := postgresql.CheckExists(db, &user, "email = ?", req.Email)
	if !exists {
		return responseData, 400, fmt.Errorf("invalid credentials")
	}

	if !utility.CompareHash(req.Password, user.Password) {
		return responseData, 400, fmt.Errorf("invalid credentials")
	}

	userData, err := user.GetUserByID(db, user.ID)
	if err != nil {
		return responseData, http.StatusInternalServerError, fmt.Errorf("unable to fetch user " + err.Error())
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

func LogoutUser(access_uuid, owner_id string, db *gorm.DB) (gin.H, int, error) {

	var (
		responseData gin.H
	)

	access_token := models.AccessToken{ID: access_uuid, OwnerID: owner_id}

	// revoke user access_token to invalidate session
	err := access_token.RevokeAccessToken(db)

	if err != nil {
		return responseData, http.StatusInternalServerError, fmt.Errorf("error revoking user session: " + err.Error())
	}

	responseData = gin.H{}

	return responseData, http.StatusOK, nil
}
