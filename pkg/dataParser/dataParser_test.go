package dataParser

import (
	"dataIngestion/pkg/models"
	apptypes "dataIngestion/types"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDataMovement(t *testing.T) {

	input := []models.Source{
		{UserID: 1, ID: 101, Title: "Test Title", Body: "Test Body"},
		{UserID: 2, ID: 102, Title: "Another", Body: "Second Body"},
	}

	expectedSource := "PlaceHolderAPI"

	post := models.Posts{
		Data:       input,
		IngestedAt: time.Now().Format(time.RFC3339),
		Source:     expectedSource,
	}

	if post.Source != expectedSource {
		t.Errorf("expected Source to be %s, got %s", expectedSource, post.Source)
	}
	if post.IngestedAt == "" {
		t.Errorf("expected IngestedAt timestamp to be set, got empty")
	}
	if len(post.Data) != 2 {
		t.Errorf("expected 2 records, got %d", len(post.Data))
	}
	if post.Data[0].Title != "Test Title" {
		t.Errorf("unexpected data in transformed result")
	}
}
func TestAPITimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
	}))
	defer server.Close()
	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	app := &apptypes.App{
		Client: client,
		Config: &apptypes.AppConfig{CydresUrl: "https://jsonplaceholder.typicode.com/posts"},
	}
	_, err := DataFetch(app)
	if err == nil {
		t.Fatal("Expected timeout error but got nil")
	}

	t.Logf("Caught expected timeout: %v", err)
}
