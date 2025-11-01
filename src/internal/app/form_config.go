package app

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/sevaho/goforms/src/config"
	"github.com/sevaho/goforms/src/internal/models"
	"github.com/sevaho/goforms/src/pkg/logger"
	"gopkg.in/yaml.v3"
)

func loadFormConfigFromFile(fc *models.FormsConfig, file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, fc); err != nil {
		return fmt.Errorf("unable to load mail template config: %w", err)
	}

	fc.Check()
	return nil
}

func loadFormConfigFromBase64String(fc *models.FormsConfig, s string) error {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return fmt.Errorf("failed to decode base64: %w", err)
	}

	if err := yaml.Unmarshal(data, fc); err == nil {
		if fc.Check() {
			return nil
		}
	}

	if err := json.Unmarshal(data, fc); err == nil {
		if fc.Check() {
			return nil
		}
	}

	return fmt.Errorf("unable to parse config as either YAML or JSON")
}

func loadFormConfig(config *config.Config) models.FormsConfig {
	var fc models.FormsConfig
	var err error

	if config.FORMS_CONFIG_FILE_PATH != "" {
		err = loadFormConfigFromFile(&fc, config.FORMS_CONFIG_FILE_PATH)
	}

	if config.FORMS_CONFIG_BASE64 != "" {
		err = loadFormConfigFromBase64String(&fc, config.FORMS_CONFIG_BASE64)
	}

	if fc.Check() {
		logger.Logger.Debug().Msgf("Loaded config: %v", fc)
		return fc
	}

	if err != nil {
		log.Fatal(err)
	}

	panic("No formconfig given")
}
