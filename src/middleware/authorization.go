package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/gusrylmubarok/mygram-backend/src/domain"
	"github.com/gusrylmubarok/mygram-backend/src/helpers"
)

func AuthorizationSocialMedia(socialMediaUseCase domain.SocialMediaUseCase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			socialMedia domain.SocialMedia
			err         error
		)

		socialMediaID := ctx.Param("socialMediaId")
		userData := ctx.MustGet("userData").(jwt.MapClaims)
		userID := string(userData["id"].(string))

		if err = socialMediaUseCase.GetByID(ctx.Request.Context(), &socialMedia, socialMediaID); err != nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, helpers.ResponseMessage{
				Status:  "fail",
				Message: fmt.Sprintf("social media with id %s doesn't exist", socialMediaID),
			})

			return
		}

		if socialMedia.UserID != userID {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, helpers.ResponseMessage{
				Status:  "unauthorized",
				Message: "you don't have permission to view or edit this social media",
			})

			return
		}
	}
}

func AuthorizationPhoto(photoUseCase domain.PhotoUseCase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			photo domain.Photo
			err   error
		)

		photoID := ctx.Param("photoId")
		userData := ctx.MustGet("userData").(jwt.MapClaims)
		userID := string(userData["id"].(string))

		if err = photoUseCase.FindById(ctx.Request.Context(), &photo, photoID); err != nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, helpers.ResponseMessage{
				Status:  "fail",
				Message: fmt.Sprintf("photo with id %s doesn't exist", photoID),
			})

			return
		}

		if photo.UserID != userID {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, helpers.ResponseMessage{
				Status:  "unauthorized",
				Message: "you don't have permission to view or edit this photo",
			})

			return
		}
	}
}

func AuthorizationComment(commentUseCase domain.CommentUseCase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			comment domain.Comment
			err     error
		)

		commentID := ctx.Param("commentId")
		userData := ctx.MustGet("userData").(jwt.MapClaims)
		userID := string(userData["id"].(string))

		if err = commentUseCase.GetByID(ctx.Request.Context(), &comment, commentID); err != nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, helpers.ResponseMessage{
				Status:  "fail",
				Message: fmt.Sprintf("comment with id %s doesn't exist", commentID),
			})

			return
		}

		if comment.UserID != userID {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, helpers.ResponseMessage{
				Status:  "unauthorized",
				Message: "you don't have permission to view or edit this comment",
			})

			return
		}
	}
}
