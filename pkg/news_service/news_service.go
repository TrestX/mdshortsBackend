package news_service

import (
	"encoding/json"
	"net/http"
	"time"

	"MdShorts/pkg/entity"
	db "MdShorts/pkg/news_service/dbs"
	"MdShorts/pkg/repository/news_repository"

	"github.com/aekam27/trestCommon"
	"github.com/gorilla/mux"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	newsService = db.NewNewsService(news_repository.NewNewsRepository("news"))
)

func GetNews(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	trestCommon.DLogMap("Get news", logrus.Fields{
		"start_time": startTime})
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// tokenString := strings.Split(r.Header.Get("Authorization"), " ")
	// if len(tokenString) < 2 {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	json.NewEncoder(w).Encode(bson.M{"status": false, "error": "authorization failed"})
	// 	return
	// }
	// _, err := trestCommon.DecodeToken(tokenString[1])
	// if err != nil {
	// 	trestCommon.ECLog1(errors.Wrapf(err, "failed to authenticate token"))
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	json.NewEncoder(w).Encode(bson.M{"status": false, "error": "authorization failed"})
	// 	return
	// }
	var userID = mux.Vars(r)["userId"]
	data, err := newsService.GetNews(userID)
	if err != nil {
		trestCommon.ECLog1(errors.Wrapf(err, "get news"))

		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "get news success"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"status": true, "error": "", "data": data})
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	trestCommon.DLogMap("get news success", logrus.Fields{
		"duration": duration,
	})
}

func GetGlobalNews(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	trestCommon.DLogMap("Get news", logrus.Fields{
		"start_time": startTime})
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	country := ""
	language := ""
	var err error
	countryS := r.URL.Query().Get("country")
	languageS := r.URL.Query().Get("language")
	if countryS != "" {
		country = countryS
	}
	if languageS != "" {
		language = languageS
	}
	data, err := newsService.GetGlobalNews(country, language)
	if err != nil {
		trestCommon.ECLog1(errors.Wrapf(err, "get news"))

		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "get news success"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"status": true, "error": "", "data": data})
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	trestCommon.DLogMap("get news success", logrus.Fields{
		"duration": duration,
	})
}

func GetSearchNews(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	trestCommon.DLogMap("Get Search News", logrus.Fields{
		"start_time": startTime})
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	search := ""
	var err error
	searchS := r.URL.Query().Get("sea")
	if searchS != "" {
		search = searchS
	}
	if searchS == "" {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bson.M{"status": true, "error": "", "data": []entity.NewsDB{}})
		endTime := time.Now()
		duration := endTime.Sub(startTime)
		trestCommon.DLogMap("get news success", logrus.Fields{
			"duration": duration,
		})
		return
	}
	data, err := newsService.GetSearchNews(search)
	if err != nil {
		trestCommon.ECLog1(errors.Wrapf(err, "get news"))

		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "get news success"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"status": true, "error": "", "data": data})
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	trestCommon.DLogMap("get news success", logrus.Fields{
		"duration": duration,
	})
}
