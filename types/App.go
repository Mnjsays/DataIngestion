package types

import (
	"context"
	"go.uber.org/zap"
)

type App struct {
	Env    string
	Ctx    context.Context
	Logger *zap.Logger
	Config *AppConfig
}
