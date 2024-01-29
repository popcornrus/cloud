package repository

import (
	"cloud/external/db"
	"cloud/internal/explorer/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

type ShareRepository struct {
	db db.MongoDBInterface
}

type ShareRepositoryInterface interface {
	Create(*model.Share) error
	FindByFileID(int64, uint64) (*model.Share, error)
	FindByUserIdAndUUID(int64, string) (*model.Share, error)
	FindByUUID(string) (*model.Share, error)
	Update(*model.Share) error
	Delete(string) error
	UpdateDownloadCount(string) error
}

func NewShareRepository(
	db db.MongoDBInterface,
) *ShareRepository {
	return &ShareRepository{
		db: db,
	}
}

func (r *ShareRepository) Create(share *model.Share) error {
	if _, err := r.db.InsertOneMongo(context.Background(), "shares", share); err != nil {
		return err
	}

	return nil
}

func (r *ShareRepository) FindByUUID(uuid string) (*model.Share, error) {
	var share model.Share

	err := r.db.FindOneMongo(context.Background(), "shares", bson.M{
		"uuid": uuid,
	}).Decode(&share)
	if err != nil {
		return nil, err
	}

	return &share, nil
}

func (r *ShareRepository) FindByUserIdAndUUID(userId int64, uuid string) (*model.Share, error) {
	var share model.Share

	err := r.db.FindOneMongo(context.Background(), "shares", bson.M{
		"uuid":    uuid,
		"user_id": userId,
	}).Decode(&share)
	if err != nil {
		return nil, err
	}

	return &share, nil
}

func (r *ShareRepository) FindByFileID(userId int64, fileID uint64) (*model.Share, error) {
	var share model.Share

	err := r.db.FindOneMongo(context.Background(), "shares", bson.M{
		"user_id": userId,
		"file_id": fileID,
	}).Decode(&share)
	if err != nil {
		return nil, err
	}

	return &share, nil
}

func (r *ShareRepository) Update(share *model.Share) error {
	_, err := r.db.UpdateOneMongo(context.Background(), "shares", bson.M{
		"_id": share.ID,
	}, bson.M{
		"$set": share,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *ShareRepository) Delete(shareUuid string) error {
	_, err := r.db.DeleteOneMongo(context.Background(), "shares", bson.M{
		"uuid": shareUuid,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *ShareRepository) UpdateDownloadCount(shareUuid string) error {
	_, err := r.db.UpdateOneMongo(context.Background(), "shares", bson.M{
		"uuid": shareUuid,
	}, bson.M{
		"$inc": bson.M{
			"download_count": 1,
		},
	})
	if err != nil {
		return err
	}

	return nil
}
