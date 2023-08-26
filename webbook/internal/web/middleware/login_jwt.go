package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-basic/webbook/internal/util"
	"net/http"
	"strings"
)

type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{
		paths: make([]string, 0),
	}
}

func (l *LoginJWTMiddlewareBuilder) Ignore(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		for _, p := range l.paths {
			if p == path {
				return
			}
		}

		headerToken := ctx.GetHeader("Authorization")
		if util.IsBlank(headerToken) {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		segs := strings.Split(headerToken, " ")
		if len(segs) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenStr := segs[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte("MxQP9pSI6BzUL9XVSZrdSeJm6Jbhw42z"), nil
		})

		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if token == nil || !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
