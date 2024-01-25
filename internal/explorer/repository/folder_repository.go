package repository

import (
	"cloud/external/db"
)

type FolderRepository struct {
	db db.MysqlInterface
}

type FolderRepositoryInterface interface {
}

func NewFolderRepository(
	db db.MysqlInterface,
) *FolderRepository {
	return &FolderRepository{
		db: db,
	}
}
