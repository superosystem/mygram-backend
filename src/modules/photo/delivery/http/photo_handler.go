package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/gusrylmubarok/mygram-backend/src/domain"
	"github.com/gusrylmubarok/mygram-backend/src/helpers"
	"github.com/gusrylmubarok/mygram-backend/src/middleware"
	"github.com/gusrylmubarok/mygram-backend/src/modules/photo/model"
)

type photoHandler struct {
	photoUseCase domain.PhotoUseCase
}

func NewPhotoHandler(routers *gin.Engine, photoUseCase domain.PhotoUseCase) {
	handler := &photoHandler{photoUseCase}

	router := routers.Group("/api/v1/photos")
	{
		router.Use(middleware.Authentication())
		router.GET("", handler.Fetch)
		router.POST("", handler.Store)
		router.PUT("/:photoId", middleware.AuthorizationPhoto(handler.photoUseCase), handler.Update)
		router.DELETE("/:photoId", middleware.AuthorizationPhoto(handler.photoUseCase), handler.Delete)
	}
}

// Fetch godoc
// @Summary    	Fetch all photos
// @Description	Get all photos with authentication user
// @Tags        photos
// @Accept      json
// @Produce     json
// @Success     200			{object}	model.ResponseDataFetchedPhoto
// @Failure     400			{object}	helpers.ResponseMessage
// @Failure     401			{object}	helpers.ResponseMessage
// @Security    Bearer
// @Router      /photos	[get]
func (handler *photoHandler) Fetch(ctx *gin.Context) {
	var (
		photos []domain.Photo
		err    error
	)

	if err = handler.photoUseCase.Fetch(ctx.Request.Context(), &photos); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	fetchedPhotos := []*model.FetchedPhoto{}

	for _, photo := range photos {
		fetchedPhotos = append(fetchedPhotos, &model.FetchedPhoto{
			ID:        photo.ID,
			Title:     photo.Title,
			Caption:   photo.Caption,
			PhotoUrl:  photo.PhotoUrl,
			UserID:    photo.UserID,
			CreatedAt: photo.CreatedAt,
			UpdatedAt: photo.UpdatedAt,
			User: &model.User{
				Email:    photo.User.Email,
				Username: photo.User.Username,
			},
		})
	}

	ctx.JSON(http.StatusOK, helpers.ResponseData{
		Status: "success",
		Data:   fetchedPhotos,
	})
}

// Store godoc
// @Summary    	Store a photo
// @Description	Create and store a photo with authentication user
// @Tags        photos
// @Accept      json
// @Produce     json
// @Param       json		body			model.AddPhoto	true	"Add Photo"
// @Success     201			{object}  		model.ResponseDataAddedPhoto
// @Failure     400			{object}		helpers.ResponseMessage
// @Failure     401			{object}		helpers.ResponseMessage
// @Security    Bearer
// @Router      /photos	[post]
func (handler *photoHandler) Store(ctx *gin.Context) {
	var (
		photo domain.Photo
		err   error
	)

	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userID := string(userData["id"].(string))

	if err = ctx.ShouldBindJSON(&photo); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	photo.UserID = userID

	if err = handler.photoUseCase.Store(ctx.Request.Context(), &photo); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusCreated, helpers.ResponseData{
		Status: "success",
		Data: model.AddedPhoto{
			ID:        photo.ID,
			Title:     photo.Title,
			Caption:   photo.Caption,
			PhotoUrl:  photo.PhotoUrl,
			UserID:    photo.UserID,
			CreatedAt: photo.CreatedAt,
		},
	})
}

// Update godoc
// @Summary     Update a photo
// @Description	Update a photo by id with authentication user
// @Tags        photos
// @Accept      json
// @Produce     json
// @Param       id		path      		string	true	"Photo ID"
// @Param       json	body			model.UpdatePhoto true  "Photo"
// @Success     200		{object}  		model.ResponseDataUpdatedPhoto
// @Failure     400		{object}		helpers.ResponseMessage
// @Failure     401		{object}		helpers.ResponseMessage
// @Failure     404		{object}		helpers.ResponseMessage
// @Security    Bearer
// @Router      /photos/{id}		[put]
func (handler *photoHandler) Update(ctx *gin.Context) {
	var (
		photo domain.Photo
		err   error
	)

	if err = ctx.ShouldBindJSON(&photo); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	updatedPhoto := domain.Photo{
		Title:    photo.Title,
		Caption:  photo.Caption,
		PhotoUrl: photo.PhotoUrl,
	}

	photoID := ctx.Param("photoId")

	if photo, err = handler.photoUseCase.Update(ctx.Request.Context(), updatedPhoto, photoID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, helpers.ResponseData{
		Status: "success",
		Data: model.UpdatedPhoto{
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
// @Summary     Delete a photo
// @Description	Delete a photo by id with authentication user
// @Tags        photos
// @Accept      json
// @Produce     json
// @Param       id	path			string	true	"Photo ID"
// @Success     200	{object}		model.ResponseMessageDeletedPhoto
// @Failure     400	{object}		helpers.ResponseMessage
// @Failure     401	{object}		helpers.ResponseMessage
// @Failure     404	{object}		helpers.ResponseMessage
// @Security    Bearer
// @Router      /photos/{id}	[delete]
func (handler *photoHandler) Delete(ctx *gin.Context) {
	photoID := ctx.Param("photoId")

	if err := handler.photoUseCase.Delete(ctx.Request.Context(), photoID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "your photo has been successfully deleted",
	})
}
