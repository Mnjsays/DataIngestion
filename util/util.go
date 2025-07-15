package util

import (
	"dataIngestion/types"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strings"
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
func Sanitize(s string) string {
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, ".", "")
	return s
}
func GetLogger() *zap.Logger {
	var Logger *zap.Logger
	logFile, err := os.OpenFile("dataingestion.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Can't open log file: %v", err)
	}

	fileWriter := zapcore.AddSync(logFile)
	consoleWriter := zapcore.AddSync(os.Stdout)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	encoder := zapcore.NewJSONEncoder(encoderConfig)

	level := zapcore.InfoLevel

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, fileWriter, level),
		zapcore.NewCore(encoder, consoleWriter, level),
	)
	Logger = zap.New(core, zap.AddCaller())
	return Logger
}
