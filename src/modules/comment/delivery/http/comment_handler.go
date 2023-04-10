package delivery

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/gusrylmubarok/mygram-backend/src/domain"
	"github.com/gusrylmubarok/mygram-backend/src/helpers"
	"github.com/gusrylmubarok/mygram-backend/src/middleware"
)

type commentHandler struct {
	commentUseCase domain.CommentUseCase
	photoUseCase   domain.PhotoUseCase
}

func NewCommentHandler(routers *gin.Engine, commentUseCase domain.CommentUseCase, photoUseCase domain.PhotoUseCase) {
	handler := &commentHandler{commentUseCase, photoUseCase}

	router := routers.Group("/api/v1/comment")
	{
		router.Use(middleware.Authentication())
		router.GET("/:photoId", handler.GetByPhoto)
		router.POST("", handler.Save)
		router.PUT("/:commentId", middleware.AuthorizationComment(handler.commentUseCase), handler.Update)
		router.DELETE("/:commentId", middleware.AuthorizationComment(handler.commentUseCase), handler.DeleteById)
	}
}

// Save godoc
// @Summary			Add a comment
// @Description		create and store a comment with authentication user
// @Tags        	comment
// @Accept      	json
// @Produce     	json
// @Param       	json	body			domain.AddComment true  "Add Comment"
// @Success     	201		{object}  		domain.ResponseAddedComment
// @Failure     	400		{object}		helpers.ResponseMessage
// @Failure     	401		{object}		helpers.ResponseMessage
// @Security    	Bearer
// @Router      	/comment	[post]
func (handler *commentHandler) Save(ctx *gin.Context) {
	var (
		input   domain.AddComment
		comment domain.Comment
		photo   domain.Photo
		err     error
	)

	if err = ctx.ShouldBindJSON(&input); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userID := string(userData["id"].(string))

	comment.UserID = userID
	comment.PhotoID = input.PhotoID
	comment.Message = input.Message

	if err = handler.photoUseCase.FindById(ctx.Request.Context(), &photo, comment.PhotoID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, helpers.ResponseMessage{
			Status:  "fail",
			Message: fmt.Sprintf("photo with id %s doesn't exist", comment.PhotoID),
		})

		return
	}

	if err = handler.commentUseCase.Save(ctx.Request.Context(), &comment); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusCreated, domain.ResponseAddedComment{
		Status: "success",
		Data: domain.AddedComment{
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
// @Tags        	comment
// @Accept      	json
// @Produce     	json
// @Param       	id		path			string  true  "Comment ID"
// @Param       	json	body			domain.UpdateComment	true	"Update Comment"
// @Success     	200		{object}  		domain.ResponseUpdatedComment
// @Failure     	400		{object}		helpers.ResponseMessage
// @Failure     	401		{object}		helpers.ResponseMessage
// @Failure     	404		{object}		helpers.ResponseMessage
// @Security    	Bearer
// @Router      	/comment/{id}	[put]
func (handler *commentHandler) Update(ctx *gin.Context) {
	var (
		input   domain.UpdateComment
		comment domain.Comment
		err     error
	)

	if err = ctx.ShouldBindJSON(&input); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})
		return
	}

	commentID := ctx.Param("commentId")
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userID := string(userData["id"].(string))

	comment.UserID = userID
	comment.Message = input.Message

	if comment, err = handler.commentUseCase.Update(ctx.Request.Context(), comment, commentID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, domain.ResponseUpdatedComment{
		Status: "success",
		Data: domain.UpdatedComment{
			ID:        comment.ID,
			Message:   comment.Message,
			UserID:    comment.UserID,
			PhotoID:   comment.PhotoID,
			UpdatedAt: comment.UpdatedAt,
		},
	})
}

// Delete godoc
// @Summary		Delete a comment
// @Description	Delete a comment by id with authentication user
// @Tags        comment
// @Accept      json
// @Produce     json
// @Param       id  path				string	true	"Comment ID"
// @Success     200 {object}			domain.ResponseMessageDeletedComment
// @Failure     400 {object}			helpers.ResponseMessage
// @Failure     401	{object}			helpers.ResponseMessage
// @Failure     404	{object}			helpers.ResponseMessage
// @Security    Bearer
// @Router      /comment/{id}	[delete]
func (handler *commentHandler) DeleteById(ctx *gin.Context) {
	commentID := ctx.Param("commentId")

	if err := handler.commentUseCase.DeleteById(ctx.Request.Context(), commentID); err != nil {
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

// Find By Photo godoc
// @Summary			Get all by photo comments
// @Description		Get all comments by photo with authentication user
// @Tags        	comment
// @Accept      	json
// @Produce     	json
// @Success     	200	{object}	domain.ResponseGetComment
// @Failure     	400	{object}	helpers.ResponseMessage
// @Failure     	401	{object}	helpers.ResponseMessage
// @Security    	Bearer
// @Router      	/comment/{photoId}     [get]
func (handler *commentHandler) GetByPhoto(ctx *gin.Context) {
	var (
		comments []domain.Comment
		err      error
	)

	photoID := ctx.Param("photoId")

	if err = handler.commentUseCase.FindByPhoto(ctx.Request.Context(), &comments, photoID); err != nil {
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

/*
// Fetch godoc
// @Summary			Fetch all comments
// @Description		Get all comments with authentication user
// @Tags        	comment
// @Accept      	json
// @Produce     	json
// @Success     	200	{object}	domain.ResponseDataFetchedComment
// @Failure     	400	{object}	helpers.ResponseMessage
// @Failure     	401	{object}	helpers.ResponseMessage
// @Security    	Bearer
// @Router      	/comment     [get]
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
*/
