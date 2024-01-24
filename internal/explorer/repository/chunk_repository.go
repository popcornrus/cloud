package repository

import (
	"cloud/internal/explorer/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type ChunkRepositoryInterface interface {
	RemoveFileChunks(context.Context, uint64) error
	GetFileChunks(context.Context, uint64) ([]*model.FileChunk, error)
	StoreFileChunk(context.Context, *model.FileChunk) error
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

func (fr *FileRepository) RemoveFileChunks(ctx context.Context, fileID uint64) error {
	_, err := fr.mongodb.DeleteManyMongo(ctx, "file_chunks", bson.D{
		{"file_id", fileID},
	})
	if err != nil {
		return err
	}

	return nil
}

func (fr *FileRepository) StoreFileChunk(ctx context.Context, fileChunk *model.FileChunk) error {
	fileChunk.CreatedAt = time.Now()

	_, err := fr.mongodb.InsertOneMongo(ctx, "file_chunks", fileChunk)
	if err != nil {
		return err
	}

	return nil
}
