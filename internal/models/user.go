package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/hngprojects/telex_be/pkg/repository/storage/postgresql"
)

type User struct {
	ID         string         `gorm:"type:uuid;primaryKey;unique;not null" json:"id"`
	Name       string         `gorm:"column:name; type:varchar(255)" json:"name"`
	Email      string         `gorm:"column:email; type:varchar(255)" json:"email"`
	IsVerified bool           `gorm:"column:is_verified; type:bool" json:"is_verified"`
	Profile    Profile        `gorm:"foreignKey:Userid;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"profile"`
	Rooms      []Room         `gorm:"many2many:user_rooms;" json:"rooms"`
	Password   string         `gorm:"column:password; type:text; not null" json:"-"`
	CreatedAt  time.Time      `gorm:"column:created_at; not null; autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at; null; autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type CreateUserRequestModel struct {
	Email       string `json:"email" validate:"required"`
	Password    string `json:"password" validate:"required"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	UserName    string `json:"username"`
	PhoneNumber string `json:"phone_number"`
}

type UpdateUserRequestModel struct {
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	UserName    string `json:"username" validate:"required"`
	PhoneNumber string `json:"phone_number"`
}

type LoginRequestModel struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (u *User) GetUserByID(db *gorm.DB, userID string) (User, error) {
	var user User

	query := db.Where("id = ?", userID)
	query = postgresql.PreloadEntities(query, &user, "Profile")

	if err := query.First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (u *User) GetUserByEmail(db *gorm.DB, userEmail string) (User, error) {
	var user User

	query := db.Where("email = ?", userEmail)
	query = postgresql.PreloadEntities(query, &user, "Profile")

	if err := query.First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (u *User) CreateUser(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &u)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Update(db *gorm.DB) error {
	_, err := postgresql.SaveAllFields(db, &u)
	return err
}

func (u *User) DeleteAUser(db *gorm.DB) error {

	err := postgresql.DeleteRecordFromDb(db, u)

	if err != nil {
		return err
	}

	return nil
}

func (u *User) GetProfileID(db *gorm.DB, userID string) (string, error) {
	var user User

	query := db.Where("id = ?", userID)
	query = postgresql.PreloadEntities(query, &user, "Profile")

	if err := query.First(&user).Error; err != nil {
		return user.Profile.ID, err
	}

	return user.Profile.ID, nil
}

func (u *User) GetUserWithProfile(db *gorm.DB, userID string) (User, error) {
	var user User

	query := db.Where("id = ?", userID)
	query = postgresql.PreloadEntities(query, &user, "Profile")

	if err := query.First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}
