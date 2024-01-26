package repository

import "cloud/external/db"

type ShareRepository struct {
	db db.MongoDBInterface
}

type ShareRepositoryInterface interface {
}

func NewShareRepository(
	db db.MongoDBInterface,
) *ShareRepository {
	return &ShareRepository{
		db: db,
	}
}
