package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gusrylmubarok/mygram-backend/src/helpers"
)

func Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		verifyToken, err := VerifyToken(ctx)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, helpers.ResponseMessage{
				Status:  "unauthenticated",
				Message: err.Error(),
			})

			return
		}

		ctx.Set("userData", verifyToken)
		ctx.Next()
	}
}
