package templates

import (
	"github.com/hngprojects/telex_be/internal/models"
	"github.com/hngprojects/telex_be/pkg/repository/storage"
	"github.com/hngprojects/telex_be/utility"
	"gorm.io/gorm"
)

func CreateTemplate(db *storage.Database, templateReq models.TemplateRequest) (models.EmailTemplate, error) {

	template := models.EmailTemplate{
		ID:   utility.GenerateUUID(),
		Name: templateReq.Name,
		Body: templateReq.Body,
	}

	err := template.Create(db.Postgresql)
	if err != nil {
		return template, err
	}
	return template, nil
}

func GetTemplates(db *gorm.DB) ([]models.EmailTemplate, error) {
	var template models.EmailTemplate

	templates, err := template.GetAllTemplates(db)
	if err != nil {
		return nil, err
	}
	return templates, nil
}

func GetTemplate(db *gorm.DB, id string) (models.EmailTemplate, error) {
	var template models.EmailTemplate

	temp, err := template.GetTemplateByID(db, id)
	if err != nil {
		return models.EmailTemplate{}, err
	}
	return temp, nil
}

func DeleteTemplate(db *gorm.DB, id string) error {
	var template models.EmailTemplate

	template.ID = id

	err := template.DeleteTemplate(db, id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateTemplate(db *gorm.DB, id string, templateReq models.EmailTemplate) (models.EmailTemplate, error) {
	var template models.EmailTemplate

	template.ID = id
	template.Name = templateReq.Name
	template.Body = templateReq.Body

	temp, err := template.UpdateTemplate(db, id)
	if err != nil {
		return temp, err
	}
	return temp, nil
}
