package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/gusrylmubarok/mygram-backend/src/domain"
	"github.com/gusrylmubarok/mygram-backend/src/helpers"
	"github.com/gusrylmubarok/mygram-backend/src/middleware"
)

type photoHandler struct {
	photoUseCase domain.PhotoUseCase
}

func NewPhotoHandler(routers *gin.Engine, photoUseCase domain.PhotoUseCase) {
	handler := &photoHandler{photoUseCase}

	router := routers.Group("/api/v1/photo")
	{
		router.Use(middleware.Authentication())
		router.POST("", handler.CreatePhoto)
		router.PUT("/:photoId", middleware.AuthorizationPhoto(handler.photoUseCase), handler.UpdatePhoto)
		router.DELETE("/:photoId", middleware.AuthorizationPhoto(handler.photoUseCase), handler.DeleteById)
		router.GET("", handler.GetAll)
		router.GET("/:photoId", handler.GetById)
	}
}

// CreatePhoto godoc
// @Summary    	Store a photo
// @Description	Create and store a photo with authentication user
// @Tags        photo
// @Accept      json
// @Produce     json
// @Param       json		body			domain.AddPhoto	true	"Add Photo"
// @Success     201			{object}  		domain.AddedPhoto
// @Failure     400			{object}		helpers.ResponseMessage
// @Failure     401			{object}		helpers.ResponseMessage
// @Security    Bearer
// @Router      /photo	[post]
func (handler *photoHandler) CreatePhoto(ctx *gin.Context) {
	var (
		input domain.AddPhoto
		photo domain.Photo
		err   error
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

	photo.Title = input.Title
	photo.Caption = input.Caption
	photo.PhotoUrl = input.PhotoUrl
	photo.UserID = userID

	if err = handler.photoUseCase.Save(ctx.Request.Context(), &photo); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, domain.AddedPhoto{
		Status:  "success",
		Message: "photo has been created",
		Data: domain.AddedDataPhoto{
			ID:       photo.ID,
			Title:    photo.Title,
			Caption:  photo.Caption,
			PhotoUrl: photo.PhotoUrl,
			User: &domain.GetUser{
				ID:       photo.User.ID,
				Email:    photo.User.Email,
				Username: photo.User.Username,
			},
			CreatedAt: photo.CreatedAt,
		},
	})
}

// Update godoc
// @Summary     Update a photo
// @Description	Update a photo by id with authentication user
// @Tags        photo
// @Accept      json
// @Produce     json
// @Param       id		path      		string	true	"Photo ID"
// @Param       json	body			domain.UpdatePhoto true  "Photo"
// @Success     200		{object}  		domain.UpdatedPhoto
// @Failure     400		{object}		helpers.ResponseMessage
// @Failure     401		{object}		helpers.ResponseMessage
// @Failure     404		{object}		helpers.ResponseMessage
// @Security    Bearer
// @Router      /photo/{id}		[put]
func (handler *photoHandler) UpdatePhoto(ctx *gin.Context) {
	var (
		input domain.UpdatePhoto
		photo domain.Photo
		err   error
	)

	if err = ctx.ShouldBindJSON(&input); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})
		return
	}

	photo.Title = input.Title
	photo.Caption = input.Caption
	photo.PhotoUrl = input.PhotoUrl

	photoID := ctx.Param("photoId")

	if photo, err = handler.photoUseCase.Update(ctx.Request.Context(), photo, photoID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, domain.UpdatedPhoto{
		Status:  "success",
		Message: "photo has been updated",
		Data: domain.UpdatedDataPhoto{
			ID:       photo.ID,
			Title:    photo.Title,
			Caption:  photo.Caption,
			PhotoUrl: photo.PhotoUrl,
			User: &domain.GetUser{
				ID:       photo.User.ID,
				Email:    photo.User.Email,
				Username: photo.User.Username,
			},
			UpdatedAt: photo.UpdatedAt,
		},
	})
}

// Delete By Id godoc
// @Summary     Delete a photo
// @Description	Delete a photo by id with authentication user
// @Tags        photo
// @Accept      json
// @Produce     json
// @Param       id	path			string	true	"Photo ID"
// @Success     200	{object}		domain.DeletedPhoto
// @Failure     400	{object}		helpers.ResponseMessage
// @Failure     401	{object}		helpers.ResponseMessage
// @Failure     404	{object}		helpers.ResponseMessage
// @Security    Bearer
// @Router      /photo/{id}	[delete]
func (handler *photoHandler) DeleteById(ctx *gin.Context) {
	photoID := ctx.Param("photoId")

	if err := handler.photoUseCase.DeleteById(ctx.Request.Context(), photoID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, domain.DeletedPhoto{
		Status:  "success",
		Message: "photo has been deleted",
	})
}

// Get All godoc
// @Summary    	Get all photos
// @Description	Get all photos
// @Tags        photo
// @Accept      json
// @Produce     json
// @Success     200			{object}	domain.GetAllPhotos
// @Failure     400			{object}	helpers.ResponseMessage
// @Failure     401			{object}	helpers.ResponseMessage
// @Security    Bearer
// @Router      /photo	[get]
func (handler *photoHandler) GetAll(ctx *gin.Context) {
	var (
		photos []domain.Photo
		err    error
	)

	if err = handler.photoUseCase.FindAll(ctx.Request.Context(), &photos); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	fetchedPhotos := []*domain.GetDetailPhoto{}

	for _, photo := range photos {
		fetchedPhotos = append(fetchedPhotos, &domain.GetDetailPhoto{
			ID:       photo.ID,
			Title:    photo.Title,
			Caption:  photo.Caption,
			PhotoUrl: photo.PhotoUrl,
			User: &domain.GetUser{
				ID:       photo.User.ID,
				Email:    photo.User.Email,
				Username: photo.User.Username,
			},
			CreatedAt: photo.CreatedAt,
			UpdatedAt: photo.UpdatedAt,
		})
	}

	ctx.JSON(http.StatusOK, domain.GetAllPhotos{
		Status:  "success",
		Message: "get all photos",
		Data:    fetchedPhotos,
	})
}

// Get One godoc
// @Summary    	Get one photo
// @Description	Get one photo with authentication user
// @Tags        photo
// @Accept      json
// @Produce     json
// @Success     200			{object}	domain.GetByIdPhoto
// @Failure     400			{object}	helpers.ResponseMessage
// @Failure     401			{object}	helpers.ResponseMessage
// @Security    Bearer
// @Router      /photo/{id}	[get]
func (handler *photoHandler) GetById(ctx *gin.Context) {
	var (
		photo domain.Photo
		err   error
	)

	photoID := ctx.Param("photoId")

	if err = handler.photoUseCase.FindById(ctx.Request.Context(), &photo, photoID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, domain.GetByIdPhoto{
		Status:  "success",
		Message: "get detail photo",
		Data: domain.GetDetailPhoto{
			ID:       photo.ID,
			Title:    photo.Title,
			Caption:  photo.Caption,
			PhotoUrl: photo.PhotoUrl,
			User: &domain.GetUser{
				ID:       photo.User.ID,
				Email:    photo.User.Email,
				Username: photo.User.Username,
			},
			CreatedAt: photo.CreatedAt,
			UpdatedAt: photo.UpdatedAt,
		},
	})
}
