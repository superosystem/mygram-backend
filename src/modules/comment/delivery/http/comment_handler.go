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
		router.POST("", handler.CreateComment)
		router.PUT("/:commentId", middleware.AuthorizationComment(handler.commentUseCase), handler.UpdateComment)
		router.DELETE("/:commentId", middleware.AuthorizationComment(handler.commentUseCase), handler.DeleteById)
		router.GET("/by-user/:userId", handler.GetAllByUser)
		router.GET("/by-photo/:photoId", handler.GetAllByPhoto)
		router.GET("/:commentId", handler.GetOne)

	}
}

// CreateComment godoc
// @Summary			Add a comment
// @Description		create and store a comment with authentication user
// @Tags        	comment
// @Accept      	json
// @Produce     	json
// @Param       	json	body			domain.AddComment true  "Add Comment"
// @Success     	201		{object}  		domain.AddedComment
// @Failure     	400		{object}		helpers.ResponseMessage
// @Failure     	401		{object}		helpers.ResponseMessage
// @Security    	Bearer
// @Router      	/comment	[post]
func (handler *commentHandler) CreateComment(ctx *gin.Context) {
	var (
		input   domain.AddComment
		comment domain.Comment
		photo   domain.Photo
		err     error
	)

	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userID := string(userData["id"].(string))

	if err = ctx.ShouldBindJSON(&input); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	if err = handler.photoUseCase.FindById(ctx.Request.Context(), &photo, input.PhotoID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, helpers.ResponseMessage{
			Status:  "fail",
			Message: fmt.Sprintf("photo with id %s doesn't exist", input.PhotoID),
		})

		return
	}

	comment.Message = input.Message
	comment.PhotoID = input.PhotoID
	comment.UserID = userID

	if err = handler.commentUseCase.Save(ctx.Request.Context(), &comment); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusCreated, domain.AddedComment{
		Status: "success",
		Data: domain.AddedDataComment{
			ID:      comment.ID,
			Message: comment.Message,
			User: &domain.GetUser{
				ID:       comment.User.ID,
				Email:    comment.User.Email,
				Username: comment.User.Username,
			},
			Photo: &domain.GetPhoto{
				ID:       comment.Photo.ID,
				Title:    comment.Photo.Title,
				Caption:  comment.Photo.Caption,
				PhotoUrl: comment.Photo.PhotoUrl,
			},
			CreatedAt: comment.CreatedAt,
		},
	})
}

// UpdateComment godoc
// @Summary			Update a comment
// @Description		Update a comment by id with authentication user
// @Tags        	comment
// @Accept      	json
// @Produce     	json
// @Param       	id		path			string  true  "Comment ID"
// @Param       	json	body			domain.UpdateComment	true	"Update Comment"
// @Success     	200		{object}  		domain.UpdatedComment
// @Failure     	400		{object}		helpers.ResponseMessage
// @Failure     	401		{object}		helpers.ResponseMessage
// @Failure     	404		{object}		helpers.ResponseMessage
// @Security    	Bearer
// @Router      	/comment/{commentId}	[put]
func (handler *commentHandler) UpdateComment(ctx *gin.Context) {
	var (
		input   domain.UpdateComment
		comment domain.Comment
		err     error
	)

	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userID := string(userData["id"].(string))

	if err = ctx.ShouldBindJSON(&input); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})
		return
	}

	comment.Message = input.Message
	comment.UserID = userID

	commentID := ctx.Param("commentId")

	if comment, err = handler.commentUseCase.Update(ctx.Request.Context(), comment, commentID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, domain.UpdatedComment{
		Status: "success",
		Data: domain.UpdatedDataComment{
			ID:      comment.ID,
			Message: comment.Message,
			User: &domain.GetUser{
				ID:       comment.User.ID,
				Email:    comment.User.Email,
				Username: comment.User.Username,
			},
			Photo: &domain.GetPhoto{
				ID:       comment.Photo.ID,
				Title:    comment.Photo.Title,
				Caption:  comment.Photo.Caption,
				PhotoUrl: comment.Photo.PhotoUrl,
			},
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
// @Param       commnetId  path				string	true	"Comment ID"
// @Success     200 {object}			domain.DeletedComment
// @Failure     400 {object}			helpers.ResponseMessage
// @Failure     401	{object}			helpers.ResponseMessage
// @Failure     404	{object}			helpers.ResponseMessage
// @Security    Bearer
// @Router      /comment/{commentId}	[delete]
func (handler *commentHandler) DeleteById(ctx *gin.Context) {
	commentID := ctx.Param("commentId")

	if err := handler.commentUseCase.DeleteById(ctx.Request.Context(), commentID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, domain.DeletedComment{
		Status:  "success",
		Message: "your comment has been successfully deleted",
	})
}

// GetAllByUser godoc
// @Summary			Get all comments
// @Description		Get all comments with authentication user
// @Tags        	comment
// @Accept      	json
// @Produce     	json
// @Success     	200	{object}	domain.GetAllComments
// @Failure     	400	{object}	helpers.ResponseMessage
// @Failure     	401	{object}	helpers.ResponseMessage
// @Security    	Bearer
// @Router      	/comment/by-user/{photoId}	[get]
func (handler *commentHandler) GetAllByUser(ctx *gin.Context) {
	var (
		comments []domain.Comment
		err      error
	)

	userID := ctx.Param("userId")

	if err = handler.commentUseCase.FindAllByUser(ctx.Request.Context(), &comments, userID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})
		return
	}

	if len(comments) < 1 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: "comments is empty",
		})
		return
	}

	ctx.JSON(http.StatusOK, domain.GetAllComments{
		Status:  "success",
		Message: "get all comments by user",
		Data:    comments,
	})
}

// Find By Photo godoc
// @Summary			Get all by photo comments
// @Description		Get all comments by photo with authentication user
// @Tags        	comment
// @Accept      	json
// @Produce     	json
// @Success     	200	{object}	domain.GetAllComments
// @Failure     	400	{object}	helpers.ResponseMessage
// @Failure     	401	{object}	helpers.ResponseMessage
// @Security    	Bearer
// @Router      	/comment/by-photo/{photoId}     [get]
func (handler *commentHandler) GetAllByPhoto(ctx *gin.Context) {
	var (
		comments []domain.Comment
		err      error
	)

	photoID := ctx.Param("photoId")

	if err = handler.commentUseCase.FindAllByPhoto(ctx.Request.Context(), &comments, photoID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})
		return
	}

	if len(comments) < 1 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: "comments is empty",
		})
		return
	}

	ctx.JSON(http.StatusOK, domain.GetAllComments{
		Status:  "success",
		Message: "get all comments by photo",
		Data:    comments,
	})
}

// Find By Id godoc
// @Summary			Get all by photo comments
// @Description		Get all comments by photo with authentication user
// @Tags        	comment
// @Accept      	json
// @Produce     	json
// @Success     	200	{object}	domain.GetAComment
// @Failure     	400	{object}	helpers.ResponseMessage
// @Failure     	401	{object}	helpers.ResponseMessage
// @Security    	Bearer
// @Router      	/comment/{commentId}     [get]
func (handler *commentHandler) GetOne(ctx *gin.Context) {
	var (
		comment domain.Comment
		err     error
	)

	commentID := ctx.Param("commentId")

	if err = handler.commentUseCase.FindById(ctx.Request.Context(), &comment, commentID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, domain.GetAComment{
		Status:  "success",
		Message: "get comment by id",
		Data:    comment,
	})
}
