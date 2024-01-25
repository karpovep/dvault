package controllers

import (
	"dvault/app"
	"dvault/constants"
	models "dvault/db/entities"
	"dvault/middleware"
	"dvault/models/dtos"
	"dvault/services"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type UserController struct {
	UserService services.IUserService
}

func NewUserController(appContext app.IAppContext) UserController {
	userService := appContext.Get(constants.UserService).(services.IUserService)
	return UserController{UserService: userService}
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	claims := middleware.GetAuthClaims(ctx)
	var user models.User
	user.UserPubId = claims.UserPubId
	err := c.UserService.CreateUser(&user)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (c *UserController) GetUser(ctx *gin.Context) {
	claims := middleware.GetAuthClaims(ctx)
	user, err := c.UserService.GetUser(claims.UserPubId)
	if err != nil {
		log.Error(err)
		if strings.Contains(err.Error(), "record not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "user has not registered yet"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (c *UserController) UpdateUser(ctx *gin.Context) {
	claims := middleware.GetAuthClaims(ctx)
	var userUpdateRequest dtos.UserUpdateRequestDto
	if err := ctx.ShouldBindJSON(&userUpdateRequest); err != nil {
		log.Debug("invalid user update request ", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if userUpdateRequest.Username == nil && userUpdateRequest.IsPublic == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "no data for the update was provided"})
		return
	}
	_, err := c.UserService.UpdateUser(claims.UserPubId, userUpdateRequest)
	if err != nil {
		log.Error("error updating user ", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (c *UserController) SearchUsers(ctx *gin.Context) {
	q := ctx.Query("q")
	if len(q) <= 3 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid query"})
		return
	}
	users, err := c.UserService.SearchUsers(q)
	if err != nil {
		log.Error("error searching users ", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"users": users})
}

func (c *UserController) RegisterUserRoutes(rg *gin.RouterGroup) {
	route := rg.Group("/user")
	route.Use(middleware.JwtAuthMiddleware())
	route.POST("/sign-in", c.GetUser)
	route.POST("/sign-up", c.CreateUser)
	route.PATCH("", c.UpdateUser)
	route.GET("/search", c.SearchUsers)
}
