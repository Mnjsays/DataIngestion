package storage

import (
	"dataIngestion/pkg/models"
)

type Storage interface {
	ReadData(fileName string) ([]byte, error)
	StoreData(posts *models.Posts) error
}
