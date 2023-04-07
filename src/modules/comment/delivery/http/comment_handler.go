package delivery

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/gusrylmubarok/mygram-backend/src/domain"
	"github.com/gusrylmubarok/mygram-backend/src/helpers"
	"github.com/gusrylmubarok/mygram-backend/src/middleware"
	"github.com/gusrylmubarok/mygram-backend/src/modules/comment/model"
)

type commentHandler struct {
	commentUseCase domain.CommentUseCase
	photoUseCase   domain.PhotoUseCase
}

func NewCommentHandler(routers *gin.Engine, commentUseCase domain.CommentUseCase, photoUseCase domain.PhotoUseCase) {
	handler := &commentHandler{commentUseCase, photoUseCase}

	router := routers.Group("/api/v1/comments")
	{
		router.Use(middleware.Authentication())
		router.GET("", handler.Fetch)
		router.POST("", handler.Store)
		router.PUT("/:commentId", middleware.AuthorizationComment(handler.commentUseCase), handler.Update)
		router.DELETE("/:commentId", middleware.AuthorizationComment(handler.commentUseCase), handler.Delete)
	}
}

// Fetch godoc
// @Summary			Fetch all comments
// @Description		Get all comments with authentication user
// @Tags        	comments
// @Accept      	json
// @Produce     	json
// @Success     	200	{object}	model.ResponseDataFetchedComment
// @Failure     	400	{object}	helpers.ResponseMessage
// @Failure     	401	{object}	helpers.ResponseMessage
// @Security    	Bearer
// @Router      	/comments     [get]
func (handler *commentHandler) Fetch(ctx *gin.Context) {
	var (
		comments []domain.Comment

		err error
	)

	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userID := string(userData["id"].(string))

	if err = handler.commentUseCase.Fetch(ctx.Request.Context(), &comments, userID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, helpers.ResponseData{
		Status: "success",
		Data:   comments,
	})
}

// Store godoc
// @Summary			Add a comment
// @Description		create and store a comment with authentication user
// @Tags        	comments
// @Accept      	json
// @Produce     	json
// @Param       	json	body			model.AddComment true  "Add Comment"
// @Success     	201		{object}  		model.ResponseDataAddedComment
// @Failure     	400		{object}		helpers.ResponseMessage
// @Failure     	401		{object}		helpers.ResponseMessage
// @Security    	Bearer
// @Router      	/comments	[post]
func (handler *commentHandler) Store(ctx *gin.Context) {
	var (
		comment domain.Comment
		photo   domain.Photo
		err     error
	)

	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userID := string(userData["id"].(string))

	if err = ctx.ShouldBindJSON(&comment); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	photoID := comment.PhotoID

	if err = handler.photoUseCase.GetByID(ctx.Request.Context(), &photo, photoID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, helpers.ResponseMessage{
			Status:  "fail",
			Message: fmt.Sprintf("photo with id %s doesn't exist", photoID),
		})

		return
	}

	comment.UserID = userID

	if err = handler.commentUseCase.Store(ctx.Request.Context(), &comment); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusCreated, helpers.ResponseData{
		Status: "success",
		Data: model.AddedComment{
			ID:        comment.ID,
			UserID:    comment.UserID,
			PhotoID:   comment.PhotoID,
			Message:   comment.Message,
			CreatedAt: comment.CreatedAt,
		},
	})
}

// Update godoc
// @Summary			Update a comment
// @Description		Update a comment by id with authentication user
// @Tags        	comments
// @Accept      	json
// @Produce     	json
// @Param       	id		path			string  true  "Comment ID"
// @Param       	json	body			model.UpdateComment	true	"Update Comment"
// @Success     	200		{object}  		model.ResponseDataUpdatedComment
// @Failure     	400		{object}		helpers.ResponseMessage
// @Failure     	401		{object}		helpers.ResponseMessage
// @Failure     	404		{object}		helpers.ResponseMessage
// @Security    	Bearer
// @Router      	/comments/{id}	[put]
func (handler *commentHandler) Update(ctx *gin.Context) {
	var (
		comment domain.Comment
		photo   domain.Photo
		err     error
	)

	commentID := ctx.Param("commentId")
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userID := string(userData["id"].(string))

	if err = ctx.ShouldBindJSON(&comment); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	updatedComment := domain.Comment{
		UserID:  userID,
		Message: comment.Message,
	}

	if photo, err = handler.commentUseCase.Update(ctx.Request.Context(), updatedComment, commentID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, helpers.ResponseData{
		Status: "success",
		Data: model.UpdatedComment{
			ID:        photo.ID,
			UserID:    photo.UserID,
			Title:     photo.Title,
			PhotoUrl:  photo.PhotoUrl,
			Caption:   photo.Caption,
			UpdatedAt: photo.UpdatedAt,
		},
	})
}

// Delete godoc
// @Summary		Delete a comment
// @Description	Delete a comment by id with authentication user
// @Tags        comments
// @Accept      json
// @Produce     json
// @Param       id  path				string	true	"Comment ID"
// @Success     200 {object}			model.ResponseMessageDeletedComment
// @Failure     400 {object}			helpers.ResponseMessage
// @Failure     401	{object}			helpers.ResponseMessage
// @Failure     404	{object}			helpers.ResponseMessage
// @Security    Bearer
// @Router      /comments/{id}	[delete]
func (handler *commentHandler) Delete(ctx *gin.Context) {
	commentID := ctx.Param("commentId")

	if err := handler.commentUseCase.Delete(ctx.Request.Context(), commentID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "your comment has been successfully deleted",
	})
}
