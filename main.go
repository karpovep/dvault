package main

import (
	"context"
	"dazer/config"
	"dazer/controllers"
	"dazer/logger"
	"dazer/services"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	cfg             *config.Config
	server          *gin.Engine
	notesService    services.NoteService
	notesController controllers.NoteController
	ctx             context.Context
	noteCollection  *mongo.Collection
	mongoClient     *mongo.Client
	err             error
)

func init() {
	cfg = config.Init("config.yml")
	logger.Init(cfg.Logger)

	// start the app
	log.Info("Starting the app...")

	ctx := context.TODO()

	mongoconn := options.Client().ApplyURI("mongodb://localhost:27017")
	mongoClient, err = mongo.Connect(ctx, mongoconn)
	if err != nil {
		log.Fatal(err)
	}
	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	log.Info("mongo connection established")

	noteCollection = mongoClient.Database("dazer").Collection("notes")
	notesService = services.NewNoteService(noteCollection, ctx)
	notesController = controllers.New(notesService)
	server = gin.Default()
}

func main() {
	defer mongoClient.Disconnect(ctx)

	basepath := server.Group("/v1")
	notesController.RegisterNoteRoutes(basepath)

	log.Fatal(server.Run(cfg.Server.Port))
}
