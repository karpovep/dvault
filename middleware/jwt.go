package middleware

import (
	"encoding/base64"
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

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearerToken := ctx.Request.Header.Get("authorization")
		if len(strings.Split(bearerToken, " ")) != 2 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "authorization header is missing"})
			ctx.Abort()
			return
		}

		token := strings.Split(bearerToken, " ")[1]
		fmt.Println(token)
		jwtParser := jwt.NewParser(jwt.WithoutClaimsValidation())

		claims := &CustomClaims{}
		parsed, parts, err := jwtParser.ParseUnverified(token, claims)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			ctx.Abort()
			return
		}
		signature, err := base64.RawURLEncoding.DecodeString(parts[2])
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			ctx.Abort()
			return
		}
		publicKeyBytes, err := base64.RawURLEncoding.DecodeString(claims.UserId)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			ctx.Abort()
			return
		}
		data, err := base64.RawURLEncoding.DecodeString(parts[1])
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			ctx.Abort()
			return
		}
		fmt.Println(string(data))
		hash := crypto.Keccak256Hash(data)

		signatureNoRecoverID := signature[:len(signature)-1] // remove recovery ID
		verified := crypto.VerifySignature(publicKeyBytes, hash.Bytes(), signatureNoRecoverID)
		if !verified {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "invalid signature!"})
			ctx.Abort()
			return
		}

		parsed.Signature = []byte(parts[2])
		ctx.Set("claims", parsed.Claims)

		ctx.Next()
	}
}