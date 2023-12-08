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
	claims := middleware.GetAuthClaims(ctx)
	note.UserID = claims.UserId
	err := nc.NoteService.UpdateNote(&note)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (nc *NoteController) DeleteNote(ctx *gin.Context) {
	noteId := ctx.Param("id")
	note, err := nc.NoteService.GetNote(noteId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	claims := middleware.GetAuthClaims(ctx)
	if note.UserID != claims.UserId {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "not allowed"})
		return
	}
	err = nc.NoteService.DeleteNote(noteId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (nc *NoteController) RegisterNoteRoutes(rg *gin.RouterGroup) {
	protected := rg.Group("/notes")
	protected.Use(middleware.JwtAuthMiddleware())
	protected.POST("/create", nc.CreateNote)
	protected.GET("/getall", nc.GetAll)
	protected.PATCH("/update", nc.UpdateNote)
	protected.DELETE("/delete/:id", nc.DeleteNote)
}
