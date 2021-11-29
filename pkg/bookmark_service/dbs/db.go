package db

import (
	"MdShorts/pkg/entity"
	bookmark "MdShorts/pkg/repository/bookmark_repository"
	news_repository "MdShorts/pkg/repository/news_repository"
	"errors"
	"time"

	"github.com/aekam27/trestCommon"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	repo = bookmark.NewBookmarkRepository("bookmark")
)
var (
	newsrepo = news_repository.NewNewsRepository("news")
)

type bookmarkService struct{}

func NewBookmarkService(repository bookmark.BookMarkRepository) BookmarkService {
	repo = repository
	return &bookmarkService{}
}
func (*bookmarkService) AddBookmark(bookmark BookMark) (string, error) {
	var bookmarkEntity entity.BookmarkDB
	if bookmark.NewsId == "" {
		return "", errors.New("newsid missing")
	}
	if bookmark.UserId == "" {
		return "", errors.New("userid missing")
	}
	bD, err := checkByNIDUID(bookmark.UserId, bookmark.NewsId)
	if err == nil {
		if bD.Status == "Active" {
			return "", errors.New("bookmark already exist")
		} else {
			id, _ := primitive.ObjectIDFromHex(bD.ID.Hex())
			setParameters := bson.M{}
			if bookmark.Status != "" {
				setParameters["status"] = "Active"
			}
			setParameters["updated_time"] = time.Now()
			filter := bson.M{"_id": id}
			set := bson.M{
				"$set": setParameters,
			}
			return repo.UpdateOne(filter, set)
		}
	}
	bookmarkEntity.ID = primitive.NewObjectID()
	bookmarkEntity.Status = "Active"
	bookmarkEntity.AddedTime = time.Now()
	bookmarkEntity.NewsId = bookmark.NewsId
	bookmarkEntity.UserId = bookmark.UserId
	return repo.InsertOne(bookmarkEntity)
}

func (*bookmarkService) UpdateBookmarkStatus(bookmark BookMark, bookmarkid string) (string, error) {
	if bookmarkid == "" {
		err := errors.New("bookmark id missing")
		trestCommon.ECLog2(
			"update bookmark",
			err,
		)
		return "", err
	}

	bD, err := checkByNIDUID(bookmark.UserId, bookmark.NewsId)
	if err != nil {
		return "", errors.New("invalid bookmark Id")
	}
	id, _ := primitive.ObjectIDFromHex(bD.ID.Hex())
	setParameters := bson.M{}
	if bookmark.Status != "" {
		setParameters["status"] = bookmark.Status
	}
	setParameters["updated_time"] = time.Now()
	filter := bson.M{"_id": id}
	set := bson.M{
		"$set": setParameters,
	}
	return repo.UpdateOne(filter, set)
}

func checkByBookmarkID(id primitive.ObjectID) (entity.BookmarkDB, error) {
	bookmark, err := repo.FindOne(bson.M{"_id": id}, bson.M{})
	if err != nil {
		trestCommon.ECLog2(
			"Get Bookmark Details section",
			err,
		)
		return bookmark, err
	}
	return bookmark, nil
}

func checkByNIDUID(id string, nid string) (entity.BookmarkDB, error) {
	bookmark, err := repo.FindOne(bson.M{"user_id": id, "newsId": nid}, bson.M{})
	if err != nil {
		trestCommon.ECLog2(
			"Get Bookmark Details section",
			err,
		)
		return bookmark, err
	}
	return bookmark, nil
}

func (*bookmarkService) GetBookmarks(limit, skip int, status, userid, newsid string) ([]entity.NewsDB, error) {
	filter := bson.M{}
	if status != "" {
		filter["status"] = status
	}
	if userid != "" {
		filter["user_id"] = userid
	}
	if newsid != "" {
		filter["news_id"] = newsid
	}
	bookmark, err := repo.Find(filter, bson.M{})
	if err != nil {
		trestCommon.ECLog2(
			"Get bookmark section",
			err,
		)
		return []entity.NewsDB{}, err
	}
	var newsList []entity.NewsDB
	for i := 0; i < len(bookmark); i++ {
		id, _ := primitive.ObjectIDFromHex(bookmark[i].NewsId)
		newdfromId, _ := newsrepo.FindOne(bson.M{"_id": id}, bson.M{})
		newsList = append(newsList, newdfromId)
	}
	return newsList, nil
}
