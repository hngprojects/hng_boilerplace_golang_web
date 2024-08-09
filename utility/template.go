package utility

import (
	"fmt"
	"os"
	"path/filepath"
)

func FindTemplateFilePath(templateName string, templateTypePath string) (string, error) {
	steps, max, found := 1, 6, false
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for steps < max {
		steps += 1
		modPath := filepath.Join(currentDir, fmt.Sprintf("services/templates%v/%v", templateTypePath, templateName))
		_, err := os.Stat(modPath)
		if err == nil || !os.IsNotExist(err) {
			currentDir, found = modPath, true
			break
		}

		pathDir := filepath.Dir(currentDir)
		currentDir = pathDir
	}

	if !found {
		return "", fmt.Errorf("template not found")
	}
	return currentDir, nil
}
