package main

import (
	"context"
	"dvault/config"
	"dvault/controllers"
	"dvault/logger"
	"dvault/services"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

var (
	cfg             *config.Config
	server          *gin.Engine
	notesService    services.NoteService
	notesController controllers.NoteController
	userService     services.UserService
	userController  controllers.UserController
	ctx             context.Context
	noteCollection  *mongo.Collection
	userCollection  *mongo.Collection
	mongoClient     *mongo.Client
	err             error
)

func init() {
	cfg = config.Init("config.yml")
	logger.Init(cfg.Logger)

	// start the app
	log.Info("Starting the app...")

	ctx := context.TODO()

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", cfg.Db.User, cfg.Db.Pass, cfg.Db.Host, cfg.Db.Port)
	mongoconn := options.Client().ApplyURI(uri)
	mongoClient, err = mongo.Connect(ctx, mongoconn)
	if err != nil {
		log.Fatal(err)
	}
	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	log.Info("mongo connection established")

	noteCollection = mongoClient.Database(cfg.Db.Name).Collection("notes")
	notesService = services.NewNoteService(noteCollection, ctx)
	notesController = controllers.NewNoteController(notesService)
	userCollection = mongoClient.Database(cfg.Db.Name).Collection("users")
	userService = services.NewUserService(userCollection, ctx)
	userController = controllers.NewUserController(userService)
	server = gin.Default()
	server.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}

func main() {
	defer mongoClient.Disconnect(ctx)

	basepath := server.Group("/v1")
	notesController.RegisterNoteRoutes(basepath)
	userController.RegisterUserRoutes(basepath)

	log.Fatal(server.Run(cfg.Server.Port))
}
