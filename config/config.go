package config

import (
	"compare/internal/models"
	"compare/pkg/logging"
	"encoding/json"
	"io"
	"os"
)

var logger = logging.GetLogger()

func GetConfig() (*models.Config, error) {
	file, err := os.Open("./config/config.json")
	if err != nil {
		logger.Errorf("error while getting Configs: %v", err)
		return nil, err
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		logger.Errorf("error while reading Config: %v", err)
		return nil, err
	}

	var cfg models.Config

	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		logger.Errorf("can not unmarshal Config: %v", err)
		return nil, err
	}

	return &cfg, nil
}
