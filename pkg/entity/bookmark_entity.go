package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookmarkDB struct {
	ID          primitive.ObjectID `bson:"_id" json:"bookmarkId"`
	UserId      string             `bson:"user_id" json:"userId,omitempty"`
	NewsId      string             `bson:"newsId" json:"newsId"`
	Status      string             `bson:"status" json:"status"`
	AddedTime   time.Time          `bson:"added_time" json:"addedTime"`
	UpdatedTime time.Time          `bson:"updated_time" json:"updatedTime"`
}
