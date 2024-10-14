package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/vadim-rm/bmstu-web-backend/internal/auth"
	"github.com/vadim-rm/bmstu-web-backend/internal/repository"
	"github.com/vadim-rm/bmstu-web-backend/pkg/logger"
	"net/http"
	"strings"
)

const jwtPrefix = "Bearer "
const tokenClaims = "tokenClaims"

type AuthMiddleware struct {
	repository repository.Token
}

func NewAuthMiddleware(repository repository.Token) *AuthMiddleware {
	return &AuthMiddleware{
		repository: repository,
	}
}

func (a *AuthMiddleware) WithAuth(ctx *gin.Context) {
	jwtStr := ctx.GetHeader("Authorization")
	if !strings.HasPrefix(jwtStr, jwtPrefix) {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	jwtStr = jwtStr[len(jwtPrefix):]

	claims, err := a.repository.GetClaims(ctx, jwtStr)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		logger.Error(err)
		return
	}

	ctx.Set(tokenClaims, claims)
}

func (a *AuthMiddleware) WithModerator(ctx *gin.Context) {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusForbidden)
		logger.Errorf("error getting claims from context: %s", err.Error())
		return
	}

	if !claims.IsModerator {
		ctx.AbortWithStatus(http.StatusForbidden)
		logger.Error("user is not moderator")
		return
	}
}
