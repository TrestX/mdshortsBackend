package main

import (
	db "MdShorts/pkg/news_service/dbs"
	"MdShorts/pkg/router"
	"log"
	"net/http"
	"time"

	"github.com/aekam27/trestCommon"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

// setupGlobalMiddleware will setup CORS
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	handleCORS := cors.AllowAll().Handler
	return handleCORS(handler)
}

// our main function
func main() {

	trestCommon.LoadConfig()
	// create router and start listen on port 8000
	router := router.NewRouter()
	go poll()
	log.Fatal(http.ListenAndServe(":6019", setupGlobalMiddleware(router)))
}
func poll() {
	_, _ = db.GetNewNews("1")
	d := time.NewTicker(7000 * time.Second)
	for {
		select {
		case du := <-d.C:
			trestCommon.DLogMap("get news success", logrus.Fields{
				"fetched_at": du,
			})
			_, _ = db.GetNewNews("1")
		}
	}
}

////md-shorts-backend.doceree.com
