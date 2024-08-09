package models

import (
	"time"

	"github.com/hngprojects/telex_be/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type EmailTemplate struct {
	ID        string         `gorm:"type:uuid;primaryKey;unique; not null" json:"id"`
	Name      string         `gorm:"column:name; type:varchar(255); not null" json:"name"`
	Body      string         `gorm:"column:body; type:text; not null" json:"body"`
	CreatedAt time.Time      `gorm:"column:created_at; not null; autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at; null; autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type TemplateRequest struct {
	Name string `json:"name" validate:"required"`
	Body string `json:"body" validate:"required"`
}

func (t *EmailTemplate) Create(db *gorm.DB) error {

	err := postgresql.CreateOneRecord(db, &t)
	if err != nil {
		return err
	}
	return nil
}

func (t *EmailTemplate) GetAllTemplates(db *gorm.DB) ([]EmailTemplate, error) {
	var templates []EmailTemplate

	err := postgresql.SelectAllFromDb(db, "", &templates, "")
	if err != nil {
		return nil, err
	}
	return templates, nil
}

func (t *EmailTemplate) GetTemplateByID(db *gorm.DB, ID string) (EmailTemplate, error) {
	var template EmailTemplate

	exists := postgresql.CheckExists(db, &template, "id = ?", ID)
	if !exists {
		return EmailTemplate{}, gorm.ErrRecordNotFound
	}

	nerr, err := postgresql.SelectOneFromDb(db, &template, "id = ?", t.ID)
	if err != nil {
		return EmailTemplate{}, nerr
	}
	return template, nil
}

func (t *EmailTemplate) DeleteTemplate(db *gorm.DB, ID string) error {

	exists := postgresql.CheckExists(db, t, "id = ?", ID)
	if !exists {
		return gorm.ErrRecordNotFound
	}

	err := postgresql.DeleteRecordFromDb(db, &t)
	if err != nil {
		return err
	}
	return nil
}

func (t *EmailTemplate) UpdateTemplate(db *gorm.DB, ID string) (EmailTemplate, error) {
	t.ID = ID

	exists := postgresql.CheckExists(db, &EmailTemplate{}, "id = ?", ID)
	if !exists {
		return EmailTemplate{}, gorm.ErrRecordNotFound
	}

	_, err := postgresql.SaveAllFields(db, t)
	if err != nil {
		return EmailTemplate{}, err
	}

	updatedTemplate := EmailTemplate{}
	err = db.First(&updatedTemplate, "id = ?", ID).Error
	if err != nil {
		return EmailTemplate{}, err
	}

	return updatedTemplate, nil
}
