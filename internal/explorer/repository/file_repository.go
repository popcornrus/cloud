package repository

import (
	"cloud/external/db"
	"cloud/internal/explorer/model"
	"context"
	"errors"
	"time"
)

type FileRepository struct {
	db      db.MysqlInterface
	mongodb db.MongoDBInterface
}

type FileRepositoryInterface interface {
	GetByUserID(context.Context, uint64) ([]*model.File, error)
	Update(context.Context, *model.File) error
	Delete(context.Context, *model.File) error
	Create(context.Context, *model.File) (uint64, error)
	FindByUUID(context.Context, string) (*model.File, error)
	FindByID(context.Context, uint64) (*model.File, error)
	Search(context.Context, uint64, string) ([]*model.File, error)

	ChunkRepositoryInterface
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

func (fr *FileRepository) FindByID(ctx context.Context, id uint64) (*model.File, error) {
	const query = "SELECT `id`, `uuid`, `name`, `path`, `hash`, `state`, `size`, `preview`, `type`, `updated_at` FROM `files` WHERE `id` = ? LIMIT 1"

	row := fr.db.GetExecer().QueryRowContext(ctx, query, id)

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

func (fr *FileRepository) Search(ctx context.Context, userID uint64, query string) ([]*model.File, error) {
	const sqlQuery = "SELECT `id`, `uuid`, `name`, `path`, `hash`, `state`, `size`, `type`, `updated_at` FROM `files` WHERE `user_id` = ? AND `name` LIKE ?"

	rows, err := fr.db.GetExecer().QueryContext(ctx, sqlQuery, userID, "%"+query+"%")
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

func (fr *FileRepository) Delete(ctx context.Context, file *model.File) error {
	const query = "DELETE FROM `files` WHERE `id` = ?"

	_, err := fr.db.GetExecer().ExecContext(ctx, query, file.ID)
	if err != nil {
		return err
	}

	return nil
}
