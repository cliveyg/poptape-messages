package main

import (
	"github.com/gorilla/mux"
	"log"
	"fmt"
	"net/http"
    "context"
    "time"
    "strconv"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	Router *mux.Router
	Mongo  *mongo.Client
    NotificationURI string
    DocumentLimit int64
}

func (a *App) Initialize(connectionString, NotificationURI, DocumentLimit string) {

	a.Router = mux.NewRouter()
	a.initializeRoutes()
    a.NotificationURI = NotificationURI
    var err error
    var iLimit int64

    iLimit, err = strconv.ParseInt(DocumentLimit, 10, 64)
    if err != nil {
        log.Fatal(err)
    }
    a.DocumentLimit = iLimit

    clientOptions := options.Client().ApplyURI(connectionString)
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    a.Mongo, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

}


func (a *App) Run(addr string) {
	log.Print(fmt.Sprintf("Server running on port [%s]", addr))
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
