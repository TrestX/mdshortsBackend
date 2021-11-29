package bookmark_repository

import (
	"MdShorts/pkg/entity"

	"go.mongodb.org/mongo-driver/bson"
)

type BookMarkRepository interface {
	InsertOne(entity.BookmarkDB) (string, error)
	FindOne(filter, projection bson.M) (entity.BookmarkDB, error)
	Find(filter, projection bson.M) ([]entity.BookmarkDB, error)
	UpdateOne(filter, update bson.M) (string, error)
	FindWithIDs(filter, projection bson.M) ([]entity.BookmarkDB, error)
	DeleteOne(filter bson.M) error
}
