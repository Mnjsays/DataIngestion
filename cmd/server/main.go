package main

import (
	"context"
	"dataIngestion/pkg/dataParser"
	"dataIngestion/pkg/storage"
	"dataIngestion/types"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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
		storageBackend := getStorage(app)
		err = storageBackend.StoreData(&posts)
		if err != nil {
			app.Logger.Error("data Ingestion failed with error", zap.Error(err))
			return
		}
		app.Logger.Info("data uploaded to cloud storage i.e., S3")

		gmux := mux.NewRouter()
		gmux.HandleFunc("/gets3Data/{filename}", dataParser.DataRetriever(app, storageBackend)).Methods("GET")
		err = http.ListenAndServe(fmt.Sprintf(":%s", app.Config.Port), gmux)
		if err != nil {
			app.Logger.Error("falied to listen at given port", zap.Error(err))
		}
	}

}
func createAppConfig() (*types.App, context.CancelFunc) {
	var err error
	err = godotenv.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading .env file: %v\n", err)
		os.Exit(1)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()
	config := &types.AppConfig{
		Port:      os.Getenv("PORT"),
		CydresUrl: os.Getenv("CYDRES_URL"),
		AwsS3: types.AwsS3Config{
			BucketName: os.Getenv("AWS_BUCKET_NAME"),
			SecretKey:  os.Getenv("AWS_SECRET_KEY"),
			AccessKey:  os.Getenv("AWS_ACCESS_KEY"),
			Region:     os.Getenv("AWS_REGION"),
			Folder:     os.Getenv("AWS_FOLDER_NAME"),
		},
	}
	//env := util.GetApplicationEnvirnoment()
	//config, err := util.GetConfig(env, logger)
	//if err != nil {
	//	logger.Fatal("failed to load config details", zap.Error(err))
	//}
	return &types.App{
		Logger: logger,
		Config: config,
		Client: &http.Client{Timeout: 5 * time.Second},
		Ctx:    ctx,
	}, cancel
}
func getStorage(app *types.App) storage.Storage {

	return &storage.S3Storage{App: app}
}
