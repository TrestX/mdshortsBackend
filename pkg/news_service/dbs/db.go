package db

import (
	"MdShorts/pkg/api"
	"MdShorts/pkg/entity"
	"MdShorts/pkg/profile_service/db"
	"MdShorts/pkg/repository"
	news_repository "MdShorts/pkg/repository/news_repository"
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/aekam27/trestCommon"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	user_news_repository "MdShorts/pkg/repository/user_news_repository"
	user_news_update_repository "MdShorts/pkg/repository/user_news_update_repository"
	uscDb "MdShorts/pkg/userNewsCheck_service/dbs"
	usuDb "MdShorts/pkg/userNewsUpdate_service/dbs"

	catdb "MdShorts/pkg/category_service/dbs"
	category_repository "MdShorts/pkg/repository/category_repository"
)

var (
	repo = news_repository.NewNewsRepository("news")
)

var (
	profileService = db.NewProfileService(repository.NewProfileRepository("users"))
)
var (
	userNewsCheckService = uscDb.NewUserNewsService(user_news_repository.NewUserNewsCheckRepository("usernewscheck"))
)

var (
	userNewsUpdateService = usuDb.NewUserNewsUpdateService(user_news_update_repository.NewUserNewsUpdateRepository("usernewsupdate"))
)

var (
	categoryService = catdb.NewCategoryService(category_repository.NewCategoryRepository("category"))
)

type newsService struct{}

func NewNewsService(repository news_repository.NewsRepository) NewsService {
	repo = repository
	return &newsService{}
}
func getTopHeadliner(news entity.TopNewsStruct) ([]entity.NewsDB, error) {
	var entityNews []entity.NewsDB
	for i := 0; i < len(news.Articles); i++ {
		if news.Articles[i].UrlToImage != "" {
			body := constructBody(news.Articles[i], "others")
			entityNews = append(entityNews, body)
		}
	}
	return entityNews, nil
}
func noUserId() ([]entity.NewsDB, error) {
	data, err := api.GetHealthTopHeadlines("us", "en", "health")
	if err != nil {
		return []entity.NewsDB{}, err
	}
	return getTopHeadliner(data)
}

func (*newsService) GetGlobalNews(language, country string) ([]entity.NewsDB, error) {
	lang := "en"
	if language != "" {
		lang = language
	}
	con := ""
	if country != "" {
		con = country
	}
	data, err := api.GetHealthTopHeadlines(con, lang, "health")
	if err != nil {
		return []entity.NewsDB{}, err
	}
	return getTopHeadliner(data)
}

func (*newsService) GetSearchNews(search string) ([]entity.NewsDB, error) {
	splitQuery := strings.Split(search, " ")
	filter := bson.M{}
	var searchList bson.A
	for i := 0; i < len(splitQuery); i++ {
		searchList = append(searchList, bson.M{"title": bson.M{"$regex": splitQuery[i], "$options": "i"}})
		searchList = append(searchList, bson.M{"description": bson.M{"$regex": splitQuery[i], "$options": "i"}})
		searchList = append(searchList, bson.M{"sourceName": bson.M{"$regex": splitQuery[i], "$options": "i"}})
		searchList = append(searchList, bson.M{"author": bson.M{"$regex": splitQuery[i], "$options": "i"}})
	}
	filter["$or"] = searchList
	return repo.FindSort(filter, bson.M{"_id": -1}, bson.M{})
}

func (*newsService) GetNews(userId string) ([]entity.NewsDB, error) {
	// if strings.Trim(userId, " ") == "" {
	// 	return noUserId()
	// }
	news, err := repo.FindSort(bson.M{}, bson.M{"_id": -1}, bson.M{})
	if err != nil {
		return []entity.NewsDB{}, err
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(news), func(i, j int) {
		news[i], news[j] = news[j], news[i]
	})
	return news, nil
	// newsIds, err := userNewsCheckService.GetUser(userId, "", "Read", 0, false)
	// if err != nil || len(newsIds) == 0 {
	// 	return repo.FindSort(bson.M{}, bson.M{"_id": -1}, bson.M{})
	// } else {
	// subFilter := bson.A{}
	// for _, id := range newsIds {
	// 	idd, _ := primitive.ObjectIDFromHex(id.NewsID)
	// 	subFilter = append(subFilter, bson.M{"_id": bson.M{"$not": bson.M{"$regex": idd.Hex(), "$options": "i"}}})
	// }
	// filter := bson.M{"$or": subFilter}

	// return repo.FindWithIDs(filter, bson.M{})
	// }
}

func GetNewNews(page string) (string, error) {
	catdata, _ := categoryService.GetCategories(100, 0, "")
	searchString := ""
	for i := 0; i < len(catdata); i++ {
		searchString = strings.Replace(catdata[i].CategoryName, " ", "%20", -1)
		news, err := api.GetNewslines(searchString, time.Now().Add(-2*time.Hour).String(), page)
		if err != nil {
			trestCommon.ECLog2(
				"error while retrieving newslines for"+searchString,
				err,
			)
		}
		_, err = insertInDB(news, catdata[i].CategoryName)
		if err != nil {
			trestCommon.ECLog2(
				"Polling Articles",
				err,
			)
		}
	}
	_, err := getHeadlinesForCountry("us", "en")
	if err != nil {
		trestCommon.ECLog2(
			"error while retrieving newslines united states",
			err,
		)
	}
	_, err = getHeadlinesForCountry("gb", "en")
	if err != nil {
		trestCommon.ECLog2(
			"error while retrieving newslines uk",
			err,
		)
	}
	_, err = getHeadlinesForCountry("in", "en")
	if err != nil {
		trestCommon.ECLog2(
			"error while retrieving newslines india",
			err,
		)
	}
	_, err = getHeadlinesForCountry("au", "en")
	if err != nil {
		trestCommon.ECLog2(
			"error while retrieving newslines australia",
			err,
		)
	}
	_, err = getHeadlinesForCountry("nz", "en")
	if err != nil {
		trestCommon.ECLog2(
			"error while retrieving newslines for newzealand",
			err,
		)
	}
	_, err = getHeadlinesForCountry("ca", "en")
	if err != nil {
		trestCommon.ECLog2(
			"error while retrieving newslines for canada",
			err,
		)
	}
	var userNewsUpdate entity.UserNewsUpdateDB
	userNewsUpdate.ID = primitive.NewObjectID()
	userNewsUpdate.LastFetched = time.Now()
	userNewsUpdate.UpdatedAt = time.Now()
	p, _ := strconv.Atoi(page)
	userNewsUpdate.CurrentPage = int64(p)
	_, err = userNewsUpdateService.AddNewsUpdateForUser(userNewsUpdate)
	return "success", nil
}

func insertInDB(news entity.TopNewsStruct, categoryName string) (string, error) {
	for i := 0; i < len(news.Articles); i++ {
		body := constructBody(news.Articles[i], categoryName)
		_, err := repo.InsertOne(body)
		if err != nil {
			trestCommon.ECLog2("quiting inserting News operation", errors.New("unable to add article => name:"+news.Articles[i].Title+", url:"+news.Articles[i].Url))
			return "", errors.New("unable to add article => url:" + news.Articles[i].Url)
		}
	}
	return "successfully added all article", nil
}

func constructBody(article entity.Article, categoryName string) entity.NewsDB {
	layout := "2006-01-02T15:04:05Z"
	var newsDB entity.NewsDB
	newsDB.Category = categoryName
	newsDB.ID = primitive.NewObjectID()
	newsDB.Author = article.Author
	newsDB.Description = article.Description
	str := article.PublishedAt
	t, _ := time.Parse(layout, str)
	newsDB.PublishedAt = t
	newsDB.SourceName = article.Source.Name
	newsDB.Title = article.Title
	newsDB.Url = article.Url
	newsDB.UrlToImage = article.UrlToImage

	newsDB.AddedTime = time.Now()
	return newsDB
}

func getHeadlinesForCountry(country, language string) (string, error) {
	data, err := api.GetHealthTopHeadlines(country, language, "health")
	if err != nil {
		return "", errors.New("error while retrieving newslines")
	}
	_, err = insertInDB(data, "others")
	if err != nil {
		trestCommon.ECLog2(
			"Getting Articles",
			err,
		)
		return "", err
	}
	return "success", nil
}
