package dataParser

import (
	"dataIngestion/pkg/models"
	"dataIngestion/pkg/storage"
	"dataIngestion/types"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func DataRetriever(app *types.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.Logger.Info("Data Retriever Api called")
		vars := mux.Vars(r) // Extract path parameters
		filename, ok := vars["filename"]
		if !ok || filename == "" {
			http.Error(w, "Filename is required", http.StatusBadRequest)
			return
		}
		fileContents, err := storage.AwsRead(filename, app)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "No such file found"})
			return
		}
		var jsonData interface{}
		if err := json.Unmarshal(fileContents, &jsonData); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to parse file contents"})
			return
		}
		app.Logger.Info("Data Fetched")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jsonData)
	}
}
func DataFetch(app *types.App) (models.Posts, error) {
	if app.Config.CydresUrl == "" {
		app.Logger.Fatal("placeholder or storage url not configured")
	}
	req, err := http.NewRequestWithContext(app.Ctx, http.MethodGet, app.Config.CydresUrl, nil)
	if err != nil {
		app.Logger.Error("Error fetching data from API", zap.Error(err))
		return models.Posts{}, err
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		app.Logger.Error("Error fetching data from API", zap.Error(err))
		return models.Posts{}, err
	}
	defer resp.Body.Close()
	var source []models.Source
	json.NewDecoder(resp.Body).Decode(&source)
	post := models.Posts{
		Data:       source,
		IngestedAt: time.Now().Format(time.RFC3339),
		Source:     "PlaceHolderAPI",
	}
	return post, nil
}
