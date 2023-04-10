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
// @Success			201		{object}		domain.RegisteredUser
// @Failure			400  	{object}		helpers.ResponseMessage
// @Failure			409  	{object}		helpers.ResponseMessage
// @Router			/user/register	[post]
func (handler *userHandler) Register(ctx *gin.Context) {
	var (
		input domain.RegisterUser
		user  domain.User
		err   error
	)

	if err = ctx.ShouldBindJSON(&input); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})
		return
	}

	user.Username = input.Username
	user.Email = input.Email
	user.Password = input.Password
	user.Age = input.Age

	if err = handler.userUseCase.Register(ctx.Request.Context(), &user); err != nil {
		if strings.Contains(err.Error(), "idx_user_username") {
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

	ctx.JSON(http.StatusCreated, domain.RegisteredUser{
		Status:  "success",
		Message: "account registration has been successful",
		Data: domain.GetUser{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
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
// @Success			200		{object}		domain.LoggedInUser
// @Failure			400		{object}		helpers.ResponseMessage
// @Failure			401		{object}		helpers.ResponseMessage
// @Router			/user/login		[post]
func (handler *userHandler) Login(ctx *gin.Context) {
	var (
		input domain.LoginUser
		user  domain.User
		err   error
		token string
	)

	if err = ctx.ShouldBindJSON(&input); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})
		return
	}

	user.Email = input.Email
	user.Password = input.Password

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

	ctx.JSON(http.StatusOK, domain.LoggedInUser{
		Status:  "success",
		Message: "user login has beed successful",
		Data: domain.Token{
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
// @Success			200			{object}  		domain.UpdatedUser
// @Failure			400			{object}		helpers.ResponseMessage
// @Failure			401			{object}		helpers.ResponseMessage
// @Failure			409			{object}		helpers.ResponseMessage
// @Security		Bearer
// @Router			/user/{userId}	[put]
func (handler *userHandler) Update(ctx *gin.Context) {
	var (
		input domain.UpdateUser
		user  domain.User
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

	user.Email = input.Email
	user.Username = input.Username
	user.Age = input.Age
	user.ID = userID

	if user, err = handler.userUseCase.Update(ctx.Request.Context(), user); err != nil {
		if strings.Contains(err.Error(), "idx_user_username") {
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

	ctx.JSON(http.StatusOK, domain.UpdatedUser{
		Status:  "success",
		Message: "update user successfully",
		Data: domain.UpdatedDataUser{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			Age:       user.Age,
			UpdatedAt: user.UpdatedAt,
		},
	})
}

// Delete godoc
// @Summary			Delete own user
// @Description		Delete own user with authentication user
// @Tags			user
// @Accept			json
// @Produce			json
// @Success			200			{object}	domain.DeletedUser
// @Failure			400			{object}	helpers.ResponseMessage
// @Failure			401			{object}	helpers.ResponseMessage
// @Failure			404			{object}	helpers.ResponseMessage
// @Security		Bearer
// @Router			/user/{userId}	[delete]
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
