package util

import (
	"dataIngestion/types"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
)

func GetApplicationEnvirnoment() string {
	env := os.Getenv("DATAENV")
	if env == "" {
		env = "local"
	}
	return env
}
func GetConfig(env string, logger *zap.Logger) (*types.AppConfig, error) {
	configs := make(map[string]*types.AppConfig)
	err := ReadConfig("pkg/config/config.yaml", configs)
	if err != nil {
		return nil, err
	}
	conf, ok := configs[env]
	if !ok {
		logger.Warn("Invalid env passed , falling back to local")
		env = "local"
		conf = configs[env]
	}
	return conf, nil
}
func ReadConfig(filePath string, config interface{}) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	err = yaml.NewDecoder(f).Decode(config)
	if err != nil {
		return err
	}
	return nil
}
