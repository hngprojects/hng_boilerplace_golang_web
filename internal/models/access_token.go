package models

import (
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
)

type AccessToken struct {
	ID                        string    `gorm:"column:id; type:uuid; not null; primaryKey; unique;" json:"id"`
	OwnerID                   string    `gorm:"column:owner_id; type:uuid; not null" json:"owner_id"`
	IsLive                    bool      `gorm:"column:is_live; type:bool; default:false; not null" json:"is_live"`
	LoginAccessToken          string    `gorm:"column:login_access_token; type:text" json:"-"`
	LoginAccessTokenExpiresIn string    `gorm:"column:login_access_token_expires_in; type:varchar(250)" json:"-"`
	CreatedAt                 time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt                 time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}

func (a *AccessToken) GetAccessTokens(db *gorm.DB) error {
	err := postgresql.SelectFirstFromDb(db, &a)
	if err != nil {
		return fmt.Errorf("token selection failed: %v", err.Error())
	}
	return nil
}

func (a *AccessToken) GetByOwnerID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &a, "owner_id = ? ", a.OwnerID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (a *AccessToken) GetByID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &a, "id = ? ", a.ID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (a *AccessToken) GetLatestByOwnerIDAndIsLive(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectLatestFromDb(db, &a, "owner_id = ? and is_live = ? ", a.OwnerID, a.IsLive)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (a *AccessToken) CreateAccessToken(db *gorm.DB, tokenData interface{}) error {
	if a.OwnerID == "" {
		return fmt.Errorf("owner id not provided to create access token")
	}

	if a.ID == "" {
		return fmt.Errorf("access id not provided to create access token")
	}

	var (
		access_token = tokenData.(map[string]string)["access_token"]
		exp          = tokenData.(map[string]string)["exp"]
	)

	a.IsLive = true
	a.LoginAccessToken = access_token
	a.LoginAccessTokenExpiresIn = exp
	err := postgresql.CreateOneRecord(db, &a)
	if err != nil {
		return fmt.Errorf("user creation failed: %v", err.Error())
	}
	return nil
}

func (a *AccessToken) RevokeAccessToken(db *gorm.DB) error {
	if a.ID == "" {
		return fmt.Errorf("access token id not provided to revoke access token")
	}
	a.IsLive = false
	_, err := postgresql.SaveAllFields(db, &a)
	return err
}
