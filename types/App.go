package types

import (
	"context"
	"go.uber.org/zap"
	"net/http"
)

type App struct {
	Env    string
	Ctx    context.Context
	Logger *zap.Logger
	Config *AppConfig
	Client *http.Client
}
