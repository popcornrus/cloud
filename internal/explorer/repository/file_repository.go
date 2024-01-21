package repository

import (
	"cloud/external/db"
	"cloud/internal/explorer/model"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type FileRepository struct {
	db      db.MysqlInterface
	mongodb db.MongoDBInterface
}

type FileRepositoryInterface interface {
	GetByUserID(context.Context, uint64) ([]*model.File, error)
	GetFileChunks(context.Context, uint64) ([]*model.FileChunk, error)
	Update(context.Context, *model.File) error
	Create(context.Context, *model.File) (uint64, error)
	StoreFileChunk(context.Context, *model.FileChunk) error
	FindByUUID(context.Context, string) (*model.File, error)
}

func NewFileRepository(
	db db.MysqlInterface,
	mongoDB db.MongoDBInterface,
) *FileRepository {
	return &FileRepository{
		db:      db,
		mongodb: mongoDB,
	}
}

func (fr *FileRepository) GetByUserID(ctx context.Context, userID uint64) ([]*model.File, error) {
	const query = "SELECT `id`, `uuid`, `name`, `path`, `hash`, `state`, `size`, `type`, `updated_at` FROM `files` WHERE `user_id` = ?"

	rows, err := fr.db.GetExecer().QueryContext(ctx, query, userID)
	if err != nil {
		return nil, errors.New("failed to get files")
	}
	defer rows.Close()

	var files []*model.File

	for rows.Next() {
		var file model.File

		err := rows.Scan(
			&file.ID,
			&file.UUID,
			&file.Name,
			&file.Path,
			&file.Hash,
			&file.State,
			&file.Size,
			&file.Type,
			&file.UpdatedAt,
		)

		if err != nil {
			return nil, errors.New("failed to scan files")
		}

		files = append(files, &file)
	}

	return files, nil
}

func (fr *FileRepository) GetFileChunks(ctx context.Context, fileID uint64) ([]*model.FileChunk, error) {
	cur, err := fr.mongodb.FindMongo(ctx, "file_chunks", bson.D{
		{"file_id", fileID},
	}, &options.FindOptions{
		Sort: bson.D{
			{"chunk", 1},
		},
	})
	if err != nil {
		return nil, err
	}

	var fileChunks []*model.FileChunk
	for cur.Next(ctx) {
		var fileChunk model.FileChunk
		err := cur.Decode(&fileChunk)
		if err != nil {
			return nil, err
		}

		fileChunks = append(fileChunks, &fileChunk)
	}

	return fileChunks, nil
}

func (fr *FileRepository) Update(ctx context.Context, file *model.File) error {
	const query = "UPDATE `files` SET `name` = ?, `path` = ?, `hash` = ?, `state` = ?, `size` = ?, `type` = ?, `preview` = ?, `updated_at` = ? WHERE `id` = ?"

	now := time.Now()
	file.UpdatedAt = &now

	_, err := fr.db.GetExecer().ExecContext(ctx, query,
		file.Name,
		file.Path,
		file.Hash,
		file.State,
		file.Size,
		file.Type,
		file.Preview,
		file.UpdatedAt,
		file.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (fr *FileRepository) Create(ctx context.Context, file *model.File) (uint64, error) {
	const query = "INSERT INTO `files` (`user_id`, `uuid`, `name`, `path`, `hash`, `state`, `size`, `type`, `created_at`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"

	file.CreatedAt = time.Now()

	result, err := fr.db.GetExecer().ExecContext(ctx, query,
		file.UserID,
		file.UUID,
		file.Name,
		file.Path,
		file.Hash,
		file.State,
		file.Size,
		file.Type,
		file.CreatedAt,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(id), nil
}

func (fr *FileRepository) StoreFileChunk(ctx context.Context, fileChunk *model.FileChunk) error {
	fileChunk.CreatedAt = time.Now()

	_, err := fr.mongodb.InsertOneMongo(ctx, "file_chunks", fileChunk)
	if err != nil {
		return err
	}

	return nil
}

func (fr *FileRepository) FindByUUID(ctx context.Context, uuid string) (*model.File, error) {
	const query = "SELECT `id`, `uuid`, `name`, `path`, `hash`, `state`, `size`, `preview`, `type`, `updated_at` FROM `files` WHERE `uuid` = ? LIMIT 1"

	row := fr.db.GetExecer().QueryRowContext(ctx, query, uuid)

	var file model.File

	err := row.Scan(
		&file.ID,
		&file.UUID,
		&file.Name,
		&file.Path,
		&file.Hash,
		&file.State,
		&file.Size,
		&file.Preview,
		&file.Type,
		&file.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &file, nil
}
