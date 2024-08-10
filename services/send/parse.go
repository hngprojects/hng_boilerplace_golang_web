package send

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/hngprojects/telex_be/external/request"
	"github.com/hngprojects/telex_be/internal/config"
	"github.com/hngprojects/telex_be/utility"
)

var (
	funcMap template.FuncMap = template.FuncMap{
		"FormatInspectionPeriod": utility.FormatInspectionPeriod,
		"numberFormat":           utility.NumberFormat,
		"add":                    utility.Add,
	}
)

func ParseTemplate(extReq request.ExternalRequest, templateFileName, baseTemplateFileName string, templateData map[string]interface{}) (string, error) {
	var (
		outputBuffer bytes.Buffer
		t            *template.Template
	)
	templateData = AddMoreMailTemplateData(extReq, templateData)

	fileName, err := utility.FindTemplateFilePath(templateFileName, "/email")
	if err != nil {
		return "", err
	}

	if baseTemplateFileName != "" {
		baseFileName, err := utility.FindTemplateFilePath(baseTemplateFileName, "/email")
		if err != nil {
			return "", err
		}

		base, err := os.ReadFile(baseFileName)
		if err != nil {
			return "", err
		}
		t = template.New("base").Funcs(funcMap)
		t, err = t.Parse(string(base))
		if err != nil {
			return "", err
		}
		t, err = t.ParseFiles(fileName)
		if err != nil {
			return "", err
		}

	} else {
		filedata, err := os.ReadFile(fileName)
		if err != nil {
			return "", errors.Wrap(err, "template not found")
		}

		t, err = template.New("email_template").Funcs(funcMap).Parse(string(filedata))
		if err != nil {
			return "", err
		}
	}

	if err2 := t.Execute(&outputBuffer, templateData); err2 != nil {
		return "", err2
	}

	body := outputBuffer.String()

	return body, nil
}

func AddMoreMailTemplateData(extReq request.ExternalRequest, data map[string]interface{}) map[string]interface{} {
	appConfig := config.GetConfig()
	accountID, ok := data["account_id"].(int)
	if !ok {
		accountIDfloat, ok := data["account_id"].(float64)
		if !ok {
			accountIDStr, ok := data["account_id"].(string)
			if ok {
				accountID, _ = (strconv.Atoi(accountIDStr))
			}
		} else {
			accountID = int(accountIDfloat)
		}
	}

	data["year"] = time.Now().Year()
	data["faq"] = appConfig.App.Url + "/faq"
	if accountID != 0 {
		data["dashboard"] = fmt.Sprintf("%v/login?account-id=%v", appConfig.App.Url, accountID)
	} else {
		data["dashboard"] = fmt.Sprintf("%v/login", appConfig.App.Url)
	}

	data["business_logo_uri"] = ""

	return data
}
