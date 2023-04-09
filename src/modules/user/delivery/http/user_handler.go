package delivery

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/gusrylmubarok/mygram-backend/src/domain"
	"github.com/gusrylmubarok/mygram-backend/src/helpers"
	"github.com/gusrylmubarok/mygram-backend/src/middleware"
)

type userHandler struct {
	userUseCase domain.UserUseCase
}

func NewUserHandler(routers *gin.Engine, userUseCase domain.UserUseCase) {
	handler := &userHandler{userUseCase}

	router := routers.Group("/api/v1/user")
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
// @Tags			user
// @Accept			json
// @Produce			json
// @Param			json	body			domain.RegisterUser true "Register User"
// @Success			201		{object}		domain.ResponseRegisteredUser
// @Failure			400  	{object}		helpers.ResponseMessage
// @Failure			409  	{object}		helpers.ResponseMessage
// @Router			/user/register	[post]
func (handler *userHandler) Register(ctx *gin.Context) {
	var (
		registerUser domain.RegisterUser
		user         domain.User
		err          error
	)

	if err = ctx.ShouldBindJSON(&registerUser); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})
		return
	}

	if user, err = handler.userUseCase.Register(ctx.Request.Context(), &registerUser); err != nil {
		if strings.Contains(err.Error(), "idx_users_username") {
			ctx.AbortWithStatusJSON(http.StatusConflict, helpers.ResponseMessage{
				Status:  "fail",
				Message: "the username you entered has been used",
			})
			return
		}

		if strings.Contains(err.Error(), "idx_user_email") {
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

	ctx.JSON(http.StatusCreated, domain.ResponseRegisteredUser{
		Status: "success",
		Data: domain.RegisteredUser{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Age:      user.Age,
		},
	})
}

// Login godoc
// @Summary			Login a user
// @Description		Authentication a user and retrieve a token
// @Tags			user
// @Accept			json
// @Produce			json
// @Param			json	body			domain.LoginUser	true	"Login User"
// @Success			200		{object}		domain.ResponseLoggedInUser
// @Failure			400		{object}		helpers.ResponseMessage
// @Failure			401		{object}		helpers.ResponseMessage
// @Router			/user/login		[post]
func (handler *userHandler) Login(ctx *gin.Context) {
	var (
		loginUser domain.LoginUser
		user      domain.User
		err       error
		token     string
	)

	if err = ctx.ShouldBindJSON(&loginUser); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	if user, err = handler.userUseCase.Login(ctx.Request.Context(), &loginUser); err != nil {
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

	ctx.JSON(http.StatusOK, domain.ResponseLoggedInUser{
		Status: "success",
		Data: domain.LoggedInUser{
			Token: token,
		},
	})
}

// Update godoc
// @Summary			Update a user
// @Description		Update a user with authentication user
// @Tags			user
// @Accept			json
// @Produce			json
// @Param			json		body			domain.UpdateUser   true  "Update User"
// @Success			200			{object}  		domain.ResponseUpdatedUser
// @Failure			400			{object}		helpers.ResponseMessage
// @Failure			401			{object}		helpers.ResponseMessage
// @Failure			409			{object}		helpers.ResponseMessage
// @Security		Bearer
// @Router			/user	[put]
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
		Data: domain.UpdatedUser{
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
// @Tags			user
// @Accept			json
// @Produce			json
// @Success			200			{object}	domain.ResponseMessageDeletedUser
// @Failure			400			{object}	helpers.ResponseMessage
// @Failure			401			{object}	helpers.ResponseMessage
// @Failure			404			{object}	helpers.ResponseMessage
// @Security		Bearer
// @Router			/user	[delete]
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
