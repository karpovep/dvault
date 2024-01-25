package db

import (
	"context"
	"dvault/app"
	"dvault/config"
	"dvault/constants"
	models "dvault/db/entities"
	"dvault/db/repositories"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init(appContext app.IAppContext) (*gorm.DB, error) {
	cfg := appContext.Get(constants.AppConfig).(*config.Config)

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.Postgres.User, cfg.Postgres.Pass, cfg.Postgres.Name, cfg.Postgres.Host, cfg.Postgres.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error("db connection error ", err.Error())
		return nil, err
	}

	if cfg.Postgres.AutoMigrate {
		// Perform database migration
		err = db.AutoMigrate(&models.User{})
		if err != nil {
			log.Fatal("migration error", err)
		}
	}

	ctx := appContext.Get(constants.Ctx).(context.Context)
	userRepository := repositories.NewUserRepository(ctx, db)
	appContext.Set(constants.UserRepository, userRepository)
	return db, nil
}
