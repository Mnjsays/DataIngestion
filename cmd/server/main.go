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
	"time"
)

func main() {
	app, cancel := createAppConfig()
	defer cancel()
	if app != nil {
		app.Logger.Info("Data ingestion pipeline configured")
		posts, err := dataParser.DataFetch(app)
		if err != nil {
			app.Logger.Error("failed to fetch data from api", zap.Error(err))
			return
		}
		app.Logger.Info("data fetched from placeholders server, Ingestion to storage in  progress")
		err = storage.AwsWrite(&posts, app)
		if err != nil {
			app.Logger.Error("data Ingestion failed with error", zap.Error(err))
			return
		}
		app.Logger.Info("data uploaded to cloud storage i.e., S3")

		gmux := mux.NewRouter()
		gmux.HandleFunc("/gets3Data/{filename}", dataParser.DataRetriever(app)).Methods("GET")
		err = http.ListenAndServe(fmt.Sprintf(":%s", app.Config.Port), gmux)
		if err != nil {
			app.Logger.Error("falied to listen at given port", zap.Error(err))
		}
	}

}
func createAppConfig() (*types.App, context.CancelFunc) {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

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
	}, cancel
}
