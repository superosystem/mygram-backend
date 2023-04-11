package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/gusrylmubarok/mygram-backend/src/domain"
	"github.com/gusrylmubarok/mygram-backend/src/helpers"
	"github.com/gusrylmubarok/mygram-backend/src/middleware"
)

type socialMediaHandler struct {
	socialMediaUseCase domain.SocialMediaUseCase
}

func NewSocialMediaHandler(routers *gin.Engine, socialMediaUseCase domain.SocialMediaUseCase) {
	handler := &socialMediaHandler{socialMediaUseCase}

	router := routers.Group("/api/v1/socialmedia")
	{
		router.Use(middleware.Authentication())
		router.POST("", handler.CreateSocialMedia)
		router.PUT("/:socialMediaId", middleware.AuthorizationSocialMedia(handler.socialMediaUseCase), handler.UpdateSocialMedia)
		router.DELETE("/:socialMediaId", middleware.AuthorizationSocialMedia(handler.socialMediaUseCase), handler.DeleteSocialMedia)
		router.GET("/by-user/:userId", handler.GetAllByUser)
		router.GET("/:socialMediaId", handler.GetBySocialMediaId)

	}
}

// CreateSocialMedia godoc
// @Summary    	Add a social media
// @Description	Create and store a social media with authentication user
// @Tags        socialmedias
// @Accept      json
// @Produce     json
// @Param       json	body			domain.AddSocialMedia true  "Add Social Media"
// @Success     201		{object}  		domain.AddedSocialMedia
// @Failure     400		{object}		helpers.ResponseMessage
// @Failure     401		{object}		helpers.ResponseMessage
// @Security    Bearer
// @Router      /socialmedias		[post]
func (handler *socialMediaHandler) CreateSocialMedia(ctx *gin.Context) {
	var (
		socialMedia domain.SocialMedia
		err         error
	)

	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userID := string(userData["id"].(string))

	if err = ctx.ShouldBindJSON(&socialMedia); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})
		return
	}

	socialMedia.UserID = userID

	if err = handler.socialMediaUseCase.Save(ctx.Request.Context(), &socialMedia); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, domain.AddedSocialMedia{
		Status: "success",
		Data: domain.AddedDataSocialMedia{
			ID:             socialMedia.ID,
			UserID:         socialMedia.UserID,
			Name:           socialMedia.Name,
			SocialMediaUrl: socialMedia.SocialMediaUrl,
			CreatedAt:      socialMedia.CreatedAt,
		},
	})
}

// UpdateSocialMedia godoc
// @Summary     Update a social media
// @Description	Update a social media by id with authentication user
// @Tags        socialmedias
// @Accept      json
// @Produce     json
// @Param       id		path      string	true	"SocialMedia ID"
// @Param		json	body				domain.UpdateSocialMedia	true	"Update Social Media"
// @Success     200		{object}			domain.UpdatedSocialMedia
// @Failure     400		{object}			helpers.ResponseMessage
// @Failure     401		{object}			helpers.ResponseMessage
// @Failure     404		{object}			helpers.ResponseMessage
// @Security    Bearer
// @Router      /socialmedias/{id} [put]
func (handler *socialMediaHandler) UpdateSocialMedia(ctx *gin.Context) {
	var (
		socialMedia domain.SocialMedia
		err         error
	)

	socialMediaID := ctx.Param("socialMediaId")
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userID := string(userData["id"].(string))

	if err = ctx.ShouldBindJSON(&socialMedia); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	updatedSocialMedia := domain.SocialMedia{
		UserID:         userID,
		Name:           socialMedia.Name,
		SocialMediaUrl: socialMedia.SocialMediaUrl,
	}

	if socialMedia, err = handler.socialMediaUseCase.Update(ctx.Request.Context(), updatedSocialMedia, socialMediaID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, domain.UpdatedSocialMedia{
		Status:  "success",
		Message: "social media has been updated",
		Data: domain.UpdatedDataSocialMedia{
			ID:             socialMedia.ID,
			Name:           socialMedia.Name,
			SocialMediaUrl: socialMedia.SocialMediaUrl,
			UserID:         socialMedia.UserID,
			UpdatedAt:      socialMedia.UpdatedAt,
		},
	})
}

// DeleteSocialMedia godoc
// @Summary     Delete a social media
// @Description	Delete a social media by id with authentication user
// @Tags        socialmedias
// @Accept      json
// @Produce     json
// @Param       id   path     	string  true  "SocialMedia ID"
// @Success     200  {object}	domain.DeletedSocialMedia
// @Failure     400  {object}	helpers.ResponseMessage
// @Failure     401  {object}	helpers.ResponseMessage
// @Failure     404  {object}	helpers.ResponseMessage
// @Security    Bearer
// @Router      /socialmedias/{id} [delete]
func (handler *socialMediaHandler) DeleteSocialMedia(ctx *gin.Context) {
	socialMediaID := ctx.Param("socialMediaId")

	if err := handler.socialMediaUseCase.DeleteById(ctx.Request.Context(), socialMediaID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, domain.DeletedSocialMedia{
		Status:  "success",
		Message: "your social media has been successfully deleted",
	})
}

// GetAllByUser godoc
// @Summary    	Fetch all social media
// @Description	Get all social media with authentication user
// @Tags        socialmedias
// @Accept      json
// @Produce     json
// @Success     200	{object}	domain.GetDataSocialMedia
// @Failure     400	{object}	helpers.ResponseMessage
// @Failure     401	{object}	helpers.ResponseMessage
// @Security    Bearer
// @Router      /socialmedia/by-user/{userId}	[get]
func (handler *socialMediaHandler) GetAllByUser(ctx *gin.Context) {
	var (
		socialMedias []domain.SocialMedia
		err          error
	)

	userID := ctx.Param("userId")

	if err = handler.socialMediaUseCase.FindAllByUser(ctx.Request.Context(), &socialMedias, userID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	if len(socialMedias) < 1 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: "social media is empty",
		})
		return
	}

	ctx.JSON(http.StatusOK, helpers.ResponseData{
		Status: "success",
		Data: domain.GetDataSocialMedia{
			SocialMedias: socialMedias,
		},
	})
}

// GetAll godoc
// @Summary    	Fetch all social media
// @Description	Get all social media with authentication user
// @Tags        socialmedias
// @Accept      json
// @Produce     json
// @Success     200	{object}	domain.ResponseDataFetchedSocialMedia
// @Failure     400	{object}	helpers.ResponseMessage
// @Failure     401	{object}	helpers.ResponseMessage
// @Security    Bearer
// @Router      /socialmedia/{socialMediaId}	[get]
func (handler *socialMediaHandler) GetBySocialMediaId(ctx *gin.Context) {
	var (
		socialMedias domain.SocialMedia
		err          error
	)

	socialMediaID := ctx.Param("socialMediaId")

	if err = handler.socialMediaUseCase.FindById(ctx.Request.Context(), &socialMedias, socialMediaID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, helpers.ResponseData{
		Status: "success",
		Data: domain.GetDataSocialMedia{
			SocialMedias: socialMedias,
		},
	})
}
