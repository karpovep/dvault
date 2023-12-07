package controllers

import (
	"dazer/middleware"
	"dazer/models"
	"dazer/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	Timestamp int
	UserId    string
}

type NoteController struct {
	NoteService services.NoteService
}

func NewNoteController(noteService services.NoteService) NoteController {
	return NoteController{NoteService: noteService}
}

func (nc *NoteController) CreateNote(ctx *gin.Context) {
	claims := middleware.GetAuthClaims(ctx)
	var note models.Note
	if err := ctx.ShouldBindJSON(&note); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	note.UserID = claims.UserId
	err := nc.NoteService.CreateNote(&note)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (nc *NoteController) GetNote(ctx *gin.Context) {
	noteId := ctx.Param("id")
	note, err := nc.NoteService.GetNote(noteId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, note)
}

func (nc *NoteController) GetAll(ctx *gin.Context) {
	claims := middleware.GetAuthClaims(ctx)
	notes, err := nc.NoteService.GetAll(claims.UserId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, notes)
}

func (nc *NoteController) UpdateNote(ctx *gin.Context) {
	var note models.Note
	if err := ctx.ShouldBindJSON(&note); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := nc.NoteService.UpdateNote(&note)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (nc *NoteController) DeleteNote(ctx *gin.Context) {
	noteId := ctx.Param("id")
	err := nc.NoteService.DeleteNote(noteId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (nc *NoteController) RegisterNoteRoutes(rg *gin.RouterGroup) {
	noteRoute := rg.Group("/note")
	noteRoute.POST("/create", nc.CreateNote)
	noteRoute.GET("/get/:id", nc.GetNote)
	noteRoute.GET("/getall", nc.GetAll)
	noteRoute.PATCH("/update", nc.UpdateNote)
	noteRoute.DELETE("/delete/:id", nc.DeleteNote)

	//supports signature verification from Vue app
	protected := rg.Group("/protected")
	protected.Use(middleware.JwtAuthMiddleware())
	protected.POST("/create", nc.CreateNote)
	protected.GET("/getall", nc.GetAll)
}
