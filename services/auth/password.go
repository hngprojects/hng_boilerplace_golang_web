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

	"github.com/hngprojects/telex_be/external/request"
	"github.com/hngprojects/telex_be/internal/config"
	"github.com/hngprojects/telex_be/internal/models"
	"github.com/hngprojects/telex_be/pkg/middleware"
	"github.com/hngprojects/telex_be/pkg/repository/storage"
	"github.com/hngprojects/telex_be/pkg/repository/storage/postgresql"
	"github.com/hngprojects/telex_be/services/actions"
	"github.com/hngprojects/telex_be/services/actions/names"
	"github.com/hngprojects/telex_be/utility"
)

func UpdateUserPassword(c *gin.Context, req models.ChangePasswordRequestModel, db *gorm.DB) (*models.User, int, error) {

	user := models.User{}

	userId, err := middleware.GetUserClaims(c, db, "user_id")
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	userID, ok := userId.(string)
	if !ok {
		return nil, http.StatusBadRequest, errors.New("user_id is not of type string")
	}

	userDataExist, err := user.GetUserByID(db, userID)
	if err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("unable to fetch user " + err.Error())
	}

	if !utility.CompareHash(req.OldPassword, userDataExist.Password) && userDataExist.Password != "" {
		return nil, http.StatusBadRequest, fmt.Errorf("old password is incorrect")
	}

	if req.OldPassword == req.NewPassword {
		return nil, http.StatusConflict, errors.New("new password cannot be the same as the old password")
	}

	hashedPassword, err := utility.HashPassword(req.NewPassword)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	userDataExist.Password = hashedPassword
	err = userDataExist.Update(db)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return &userDataExist, http.StatusOK, nil
}

func PasswordReset(userEmail string, db *gorm.DB, extReq request.ExternalRequest) (string, int, error) {

	var (
		user      = models.User{}
		passReset = models.PasswordReset{}
		config    = config.GetConfig()
	)

	resetExist, err := passReset.GetPasswordResetByEmail(db, userEmail)
	if err != nil {
		return "error", http.StatusUnauthorized, err
	}

	if resetExist != nil {
		if err := resetExist.DeletePasswordReset(db); err != nil {
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

	reset := models.PasswordReset{
		ID:        utility.GenerateUUID(),
		Email:     strings.ToLower(userEmail),
		Token:     strconv.Itoa(resetToken),
		ExpiresAt: time.Now().Add(time.Duration(config.App.ResetPasswordDuration) * time.Minute),
	}

	err = reset.CreatePasswordReset(db)
	if err != nil {
		return "error", http.StatusInternalServerError, err
	}

	resetReq := models.SendOTP{
		Email:    userEmail,
		OtpToken: resetToken,
	}

	err = actions.AddNotificationToQueue(storage.DB.Redis, names.SendOTP, resetReq)
	if err != nil {
		return "error", http.StatusInternalServerError, err
	}

	return "success", http.StatusOK, nil
}

func VerifyPasswordResetToken(req models.ResetPasswordRequestModel, db *gorm.DB) (*models.User, int, error) {

	var (
		user      = models.User{}
		passReset = models.PasswordReset{}
	)

	resetExist, err := passReset.GetPasswordResetByToken(db, req.Token)
	if err != nil {
		return nil, http.StatusUnauthorized, errors.New("invalid or expired token")
	}

	userDataExist, err := user.GetUserByEmail(db, resetExist.Email)
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	hashedPassword, err := utility.HashPassword(req.NewPassword)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	userDataExist.Password = hashedPassword
	err = userDataExist.Update(db)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if err := resetExist.DeletePasswordReset(db); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &userDataExist, http.StatusOK, nil

}
