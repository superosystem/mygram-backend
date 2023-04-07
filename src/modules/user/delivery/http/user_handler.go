package delivery

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/gusrylmubarok/mygram-backend/src/domain"
	"github.com/gusrylmubarok/mygram-backend/src/helpers"
	"github.com/gusrylmubarok/mygram-backend/src/middleware"
	"github.com/gusrylmubarok/mygram-backend/src/modules/user/model"
)

type userHandler struct {
	userUseCase domain.UserUseCase
}

func NewUserHandler(routers *gin.Engine, userUseCase domain.UserUseCase) {
	handler := &userHandler{userUseCase}

	router := routers.Group("/api/v1/users")
	{
		router.POST("/register", handler.Register)
		router.POST("/login", handler.Login)
		router.PUT("", middleware.Authentication(), handler.Update)
		router.DELETE("", middleware.Authentication(), handler.Delete)
	}
}

// Register godoc
// @Summary			Register a user
// @Description		create and store a user
// @Tags			users
// @Accept			json
// @Produce			json
// @Param			json	body			model.RegisterUser	true	"Register User"
// @Success			201		{object}		model.ResponseDataRegisteredUser
// @Failure			400  	{object}		helpers.ResponseMessage
// @Failure			409  	{object}		helpers.ResponseMessage
// @Router			/users/register	[post]
func (handler *userHandler) Register(ctx *gin.Context) {
	var (
		user domain.User
		err  error
	)

	if err = ctx.ShouldBindJSON(&user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	if err = handler.userUseCase.Register(ctx.Request.Context(), &user); err != nil {
		if strings.Contains(err.Error(), "idx_users_username") {
			ctx.AbortWithStatusJSON(http.StatusConflict, helpers.ResponseMessage{
				Status:  "fail",
				Message: "the username you entered has been used",
			})

			return
		}

		if strings.Contains(err.Error(), "idx_users_email") {
			ctx.AbortWithStatusJSON(http.StatusConflict, helpers.ResponseMessage{
				Status:  "fail",
				Message: "the email you entered has been used",
			})

			return
		}

		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusCreated, helpers.ResponseData{
		Status: "success",
		Data: model.RegisteredUser{
			Age:      user.Age,
			Email:    user.Email,
			ID:       user.ID,
			Username: user.Username,
		},
	})
}

// Login godoc
// @Summary			Login a user
// @Description		Authentication a user and retrieve a token
// @Tags			users
// @Accept			json
// @Produce			json
// @Param			json	body			model.LoginUser	true	"Login User"
// @Success			200		{object}		model.ResponseDataLoggedinUser
// @Failure			400		{object}		helpers.ResponseMessage
// @Failure			401		{object}		helpers.ResponseMessage
// @Router			/users/login		[post]
func (handler *userHandler) Login(ctx *gin.Context) {
	var (
		user  domain.User
		err   error
		token string
	)

	if err = ctx.ShouldBindJSON(&user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	if err = handler.userUseCase.Login(ctx.Request.Context(), &user); err != nil {
		if strings.Contains(err.Error(), "the credential you entered are wrong") {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
				Status:  "unauthenticated",
				Message: err.Error(),
			})

			return
		}

		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "unauthenticated",
			Message: err.Error(),
		})

		return
	}

	if token = middleware.GenerateToken(user.ID, user.Email); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "unauthenticated",
			"message": err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, helpers.ResponseData{
		Status: "success",
		Data: model.LoggedinUser{
			Token: token,
		},
	})
}

// Update godoc
// @Summary			Update a user
// @Description		Update a user with authentication user
// @Tags			users
// @Accept			json
// @Produce			json
// @Param			json		body			model.UpdateUser   true  "Update User"
// @Success			200			{object}  		model.ResponseDataUpdatedUser
// @Failure			400			{object}		helpers.ResponseMessage
// @Failure			401			{object}		helpers.ResponseMessage
// @Failure			409			{object}		helpers.ResponseMessage
// @Security		Bearer
// @Router			/users	[put]
func (handler *userHandler) Update(ctx *gin.Context) {
	var (
		user domain.User
		err  error
	)

	userData := ctx.MustGet("userData").(jwt.MapClaims)
	_ = string(userData["id"].(string))

	if err = ctx.ShouldBindJSON(&user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	updatedUser := domain.User{
		Username: user.Username,
		Email:    user.Email,
	}

	if user, err = handler.userUseCase.Update(ctx.Request.Context(), updatedUser); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, helpers.ResponseData{
		Status: "success",
		Data: model.UpdatedUser{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			Age:       user.Age,
			UpdatedAt: user.UpdatedAt,
		},
	})
}

// Delete godoc
// @Summary			Delete a user
// @Description		Delete a user with authentication user
// @Tags			users
// @Accept			json
// @Produce			json
// @Success			200			{object}	model.ResponseMessageDeletedUser
// @Failure			400			{object}	helpers.ResponseMessage
// @Failure			401			{object}	helpers.ResponseMessage
// @Failure			404			{object}	helpers.ResponseMessage
// @Security		Bearer
// @Router			/users	[delete]
func (handler *userHandler) Delete(ctx *gin.Context) {
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userID := string(userData["id"].(string))

	if err := handler.userUseCase.Delete(ctx, userID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: "account not found",
		})

		return
	}

	ctx.JSON(
		http.StatusOK,
		helpers.ResponseMessage{
			Status:  "success",
			Message: "your account has been successfully deleted",
		},
	)
}
