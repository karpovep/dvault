package controllers

import (
	"dazer/middleware"
	"dazer/models"
	"dazer/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type UserController struct {
	UserService services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return UserController{UserService: userService}
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	claims := middleware.GetAuthClaims(ctx)
	var user models.User
	user.UserID = claims.UserId
	err := c.UserService.CreateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (c *UserController) GetUser(ctx *gin.Context) {
	claims := middleware.GetAuthClaims(ctx)
	user, err := c.UserService.GetUser(claims.UserId)
	fmt.Println(user)
	fmt.Println(err)
	if err != nil {
		if strings.Contains(err.Error(), "no documents") {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "user has not registered yet"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (c *UserController) RegisterUserRoutes(rg *gin.RouterGroup) {
	route := rg.Group("/user")
	route.Use(middleware.JwtAuthMiddleware())
	route.POST("/sign-in", c.GetUser)
	route.POST("/sign-up", c.CreateUser)
}
