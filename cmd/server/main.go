package main

import (
	"context"
	"dataIngestion/pkg/dataParser"
	"dataIngestion/pkg/storage"
	"dataIngestion/types"
	"dataIngestion/util"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"os"
)

func main() {
	app := createAppConfig()
	if app != nil {
		app.Logger.Info("Data ingestion pipeline configured")
		posts, err := dataParser.DataFetch(app)
		if err != nil {
			app.Logger.Error("failed to fetch data from api", zap.Error(err))
			return
		}
		app.Logger.Info("data fetched from placeholders server, Ingestion to storage in  progress")
		err = storage.AwsStorage(&posts, app)
		if err != nil {
			app.Logger.Error("data Ingestion failed with error", zap.Error(err))
			return
		}
		app.Logger.Info("data uploaded to cloud storage i.e., S3")

		gmux := mux.NewRouter()
		gmux.HandleFunc("getData", dataParser.DataRetriver).Methods("GET")
		err = http.ListenAndServe(fmt.Sprintf(":%s", app.Config.Port), gmux)
		if err != nil {
			app.Logger.Error("falied to listen at given port", zap.Error(err))
		}
	}

}
func createAppConfig() *types.App {
	var err error
	ctx := context.Background()
	env := util.GetApplicationEnvirnoment()
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	config, err := util.GetConfig(env, logger)
	if err != nil {
		logger.Fatal("failed to load config details", zap.Error(err))
	}
	return &types.App{
		Logger: logger,
		Config: config,
		Env:    env,
		Ctx:    ctx,
	}
}
