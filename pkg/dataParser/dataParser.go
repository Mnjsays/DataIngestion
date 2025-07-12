package dataParser

import (
	"dataIngestion/pkg/models"
	"dataIngestion/types"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func DataRetriver(w http.ResponseWriter, r *http.Request) {

}
func DataTransfer(app *types.App) error {
	if app.Config.CydresUrl == "" || app.Config.StorageUrl == " " {
		app.Logger.Fatal("placeholder or storage url not configured")
	}
	get, err := http.Get(app.Config.CydresUrl)
	if err != nil {
		app.Logger.Error("Error fetching data from API", zap.Error(err))
		return err
	}
	var source []models.Source
	json.NewDecoder(get.Body).Decode(&source)
	post := models.Posts{
		Data:       source,
		IngestedAt: time.Now().Format(time.RFC3339),
		Source:     "PlaceHolderAPI",
	}
	_ = post
	jsonData, err := json.MarshalIndent(post, "", " ")
	fmt.Println(string(jsonData))

	return nil
}
