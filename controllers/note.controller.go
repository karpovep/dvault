package controllers

import (
	"dazer/middleware"
	"dazer/models"
	"dazer/services"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	Timestamp int
	UserId    string
}

type NoteController struct {
	NoteService services.NoteService
}

func New(noteService services.NoteService) NoteController {
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

func (nc *NoteController) ProtectedTest(ctx *gin.Context) {
	bearerToken := ctx.Request.Header.Get("authorization")
	if len(strings.Split(bearerToken, " ")) != 2 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "authorization header is missing"})
		return
	}

	token := strings.Split(bearerToken, " ")[1]
	fmt.Println(token)
	jwtParser := jwt.NewParser(jwt.WithoutClaimsValidation())

	claims := &CustomClaims{}
	parsed, parts, err := jwtParser.ParseUnverified(token, claims)

	//parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
	//	fmt.Println(token)
	//	return []byte("test"), nil
	//})
	fmt.Println(parsed.Claims)
	fmt.Println(parsed.Method)
	fmt.Println(parsed.Header)
	fmt.Println(parsed.Signature)
	fmt.Println(parts)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	publicKeyBytes, err := base64.RawURLEncoding.DecodeString(claims.UserId)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	data, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	fmt.Println(string(data))
	hash := crypto.Keccak256Hash(data)

	signatureNoRecoverID := signature[:len(signature)-1] // remove recovery ID
	verified := crypto.VerifySignature(publicKeyBytes, hash.Bytes(), signatureNoRecoverID)

	parsed.Signature = []byte(parts[2])
	ctx.JSON(http.StatusOK, gin.H{
		"header":    parsed.Header,
		"payload":   parsed.Claims,
		"signature": parts[2],
		"verified":  verified,
	})
}

func (nc *NoteController) ProtectedTestNew(ctx *gin.Context) {
	bearerToken := ctx.Request.Header.Get("authorization")
	if len(strings.Split(bearerToken, " ")) != 2 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "authorization header is missing"})
		return
	}

	token := strings.Split(bearerToken, " ")[1]
	fmt.Println(token)
	jwtParser := jwt.NewParser(jwt.WithoutClaimsValidation())

	claims := &CustomClaims{}
	parsed, parts, err := jwtParser.ParseUnverified(token, claims)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	fmt.Println(parsed.Claims)
	fmt.Println(parsed.Method)
	fmt.Println(parsed.Header)
	fmt.Println(parsed.Signature)
	fmt.Println(parts)
	fmt.Println(parts[2])

	fmt.Println("signature", parts[2])
	signature, err := hex.DecodeString(parts[2])
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	data, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	fmt.Println(string(data))
	hash := crypto.Keccak256Hash(data)

	publicKeyBytes, err := hex.DecodeString(claims.UserId)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	signatureNoRecoverID := signature[:len(signature)-1] // remove recovery ID
	verified := crypto.VerifySignature(publicKeyBytes, hash.Bytes(), signatureNoRecoverID)

	parsed.Signature = []byte(parts[2])
	ctx.JSON(http.StatusOK, gin.H{
		"header":    parsed.Header,
		"payload":   parsed.Claims,
		"signature": parts[2],
		"verified":  verified,
	})
}

func (nc *NoteController) ProtectedTestMiddleware(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	fmt.Println("ProtectedTestMiddleware", claims)
	ctx.JSON(http.StatusOK, gin.H{
		"payload": claims,
	})
}

func (nc *NoteController) RegisterNoteRoutes(rg *gin.RouterGroup) {
	noteRoute := rg.Group("/note")
	noteRoute.POST("/create", nc.CreateNote)
	noteRoute.GET("/get/:id", nc.GetNote)
	noteRoute.GET("/getall", nc.GetAll)
	noteRoute.PATCH("/update", nc.UpdateNote)
	noteRoute.DELETE("/delete/:id", nc.DeleteNote)

	protected := rg.Group("/protected")
	protected.GET("/test", nc.ProtectedTest)
	protected.GET("/test-new", nc.ProtectedTestNew) //supports signature verification from Vue app
	protected.Use(middleware.JwtAuthMiddleware())
	protected.POST("/create", nc.CreateNote)
	protected.GET("/getall", nc.GetAll)
	protected.GET("/test-middleware", nc.ProtectedTestMiddleware)
}
