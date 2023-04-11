package delivery_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gusrylmubarok/mygram-backend/src/domain"
	mocksUseCase "github.com/gusrylmubarok/mygram-backend/src/domain/mocks/usecase"
	delivery "github.com/gusrylmubarok/mygram-backend/src/modules/user/delivery/http"
	userUseCase "github.com/gusrylmubarok/mygram-backend/src/modules/user/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUserUseCase := new(mocksUseCase.UserUseCase)
	userUseCase := userUseCase.NewUserUseCase(mockUserUseCase)

	t.Run("should success register user", func(t *testing.T) {
		// prepare
		tempMockRegisterUser := domain.User{
			Email:    "johndoe@example.com",
			Username: "johndoe",
			Password: "secret",
			Age:      8,
		}
		expected := domain.RegisteredUser{
			Status:  "success",
			Message: "user registration has been successful",
			Data: domain.GetUser{
				ID:       "",
				Username: "johndoe",
				Email:    tempMockRegisterUser.Email,
			},
		}

		mockUserUseCase.On("Register", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil).Once()

		router := gin.Default()
		rec := httptest.NewRecorder()

		userHandler := delivery.NewUserHandler(router, userUseCase)
		router.POST("/user/register", userHandler.Register)

		// do
		reqBody, err := json.Marshal(tempMockRegisterUser)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/user/register", strings.NewReader(string(reqBody)))
		router.ServeHTTP(rec, req)

		// struct to json
		// exp, err := json.Marshal(tempMockRegisteredUser)
		// assert.NoError(t, err)
		// convert to struct
		var res domain.RegisteredUser
		err = json.Unmarshal(rec.Body.Bytes(), &res)
		assert.NoError(t, err)

		// assert
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, expected.Data.Email, res.Data.Email)
	})

	t.Run("should fail register user with empty email", func(t *testing.T) {
		// prepare
		tempMockRegisterUser := domain.User{
			Email:    "",
			Username: "johndoe",
			Password: "secret",
			Age:      8,
		}

		mockUserUseCase.On("Register", mock.Anything, mock.AnythingOfType("*domain.User")).Return(errors.New("email: non zero value required")).Once()

		router := gin.Default()
		rec := httptest.NewRecorder()

		userHandler := delivery.NewUserHandler(router, userUseCase)
		router.POST("/user/register", userHandler.Register)

		// do
		reqBody, err := json.Marshal(tempMockRegisterUser)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/user/register", strings.NewReader(string(reqBody)))
		router.ServeHTTP(rec, req)

		var res domain.RegisteredUser
		err = json.Unmarshal(rec.Body.Bytes(), &res)
		assert.NoError(t, err)

		_, err = govalidator.ValidateStruct(tempMockRegisterUser)

		// assert
		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "email: non zero value required", res.Message)
	})

	t.Run("should fail register user with used email", func(t *testing.T) {
		// prepare
		tempMockRegisterUser := domain.User{
			Email:    "",
			Username: "johndoe",
			Password: "secret",
			Age:      8,
		}

		mockUserUseCase.On("Register", mock.Anything, mock.AnythingOfType("*domain.User")).Return(errors.New("duplicate key on idx_users_email")).Once()

		router := gin.Default()
		rec := httptest.NewRecorder()

		userHandler := delivery.NewUserHandler(router, userUseCase)
		router.POST("/user/register", userHandler.Register)

		// do
		reqBody, err := json.Marshal(tempMockRegisterUser)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/user/register", strings.NewReader(string(reqBody)))
		router.ServeHTTP(rec, req)

		var res domain.RegisteredUser
		err = json.Unmarshal(rec.Body.Bytes(), &res)
		assert.NoError(t, err)

		_, err = govalidator.ValidateStruct(tempMockRegisterUser)

		// assert
		assert.Error(t, err)
		assert.Equal(t, http.StatusConflict, rec.Code)
		assert.Equal(t, "the email you entered has been used", res.Message)
	})

	t.Run("should fail register user with invalid email", func(t *testing.T) {
		// prepare
		tempMockRegisterUser := domain.User{
			Email:    "johndoe.com",
			Username: "johndoe",
			Password: "secret",
			Age:      8,
		}

		mockUserUseCase.On("Register", mock.Anything, mock.AnythingOfType("*domain.User")).Return(errors.New("email: jhondoe.com does not validate as email")).Once()

		router := gin.Default()
		rec := httptest.NewRecorder()

		userHandler := delivery.NewUserHandler(router, userUseCase)
		router.POST("/user/register", userHandler.Register)

		// do
		reqBody, err := json.Marshal(tempMockRegisterUser)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/user/register", strings.NewReader(string(reqBody)))
		router.ServeHTTP(rec, req)

		var res domain.RegisteredUser
		err = json.Unmarshal(rec.Body.Bytes(), &res)
		assert.NoError(t, err)

		_, err = govalidator.ValidateStruct(tempMockRegisterUser)

		// assert
		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "email: jhondoe.com does not validate as email", res.Message)
	})

	t.Run("should fail register user with empty username", func(t *testing.T) {
		// prepare
		tempMockRegisterUser := domain.User{
			Email:    "johndoe@example.com",
			Username: "",
			Password: "secret",
			Age:      8,
		}

		mockUserUseCase.On("Register", mock.Anything, mock.AnythingOfType("*domain.User")).Return(errors.New("username: non zero value required")).Once()

		router := gin.Default()
		rec := httptest.NewRecorder()

		userHandler := delivery.NewUserHandler(router, userUseCase)
		router.POST("/user/register", userHandler.Register)

		// do
		reqBody, err := json.Marshal(tempMockRegisterUser)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/user/register", strings.NewReader(string(reqBody)))
		router.ServeHTTP(rec, req)

		var res domain.RegisteredUser
		err = json.Unmarshal(rec.Body.Bytes(), &res)
		assert.NoError(t, err)

		_, err = govalidator.ValidateStruct(tempMockRegisterUser)

		// assert
		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "username: non zero value required", res.Message)
	})

	t.Run("should fail register user with used username", func(t *testing.T) {
		// prepare
		tempMockRegisterUser := domain.User{
			Email:    "johndoe@example.com",
			Username: "johndoe",
			Password: "secret",
			Age:      8,
		}

		mockUserUseCase.On("Register", mock.Anything, mock.AnythingOfType("*domain.User")).Return(errors.New("duplicate key on idx_users_username")).Once()

		router := gin.Default()
		rec := httptest.NewRecorder()

		userHandler := delivery.NewUserHandler(router, userUseCase)
		router.POST("/user/register", userHandler.Register)

		// do
		reqBody, err := json.Marshal(tempMockRegisterUser)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/user/register", strings.NewReader(string(reqBody)))
		router.ServeHTTP(rec, req)

		var res domain.RegisteredUser
		err = json.Unmarshal(rec.Body.Bytes(), &res)
		assert.NoError(t, err)

		// assert
		assert.Equal(t, http.StatusConflict, rec.Code)
		assert.Equal(t, "the username you entered has been used", res.Message)
	})

	t.Run("should fail register user with invalid password", func(t *testing.T) {
		// prepare
		tempMockRegisterUser := domain.User{
			ID:       "user-123",
			Email:    "johndoe@example.com",
			Username: "johndoe",
			Password: "sec",
			Age:      5,
		}

		mockUserUseCase.On("Register", mock.Anything, mock.AnythingOfType("*domain.User")).Return(errors.New("password: pass does not validate as minstringlength(6)")).Once()

		router := gin.Default()
		rec := httptest.NewRecorder()

		userHandler := delivery.NewUserHandler(router, userUseCase)
		router.POST("/user/register", userHandler.Register)

		// do
		reqBody, err := json.Marshal(tempMockRegisterUser)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/user/register", strings.NewReader(string(reqBody)))
		router.ServeHTTP(rec, req)

		var res domain.RegisteredUser
		err = json.Unmarshal(rec.Body.Bytes(), &res)
		assert.NoError(t, err)

		_, err = govalidator.ValidateStruct(tempMockRegisterUser)

		// assert
		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "password: pass does not validate as minstringlength(6)", res.Message)
	})

	t.Run("should fail register user with under age", func(t *testing.T) {
		// prepare
		tempMockRegisterUser := domain.User{
			ID:       "user-123",
			Email:    "johndoe@example.com",
			Username: "johndoe",
			Password: "secret",
			Age:      5,
		}

		mockUserUseCase.On("Register", mock.Anything, mock.AnythingOfType("*domain.User")).Return(errors.New("age: 5 does not validate as range(8|63)")).Once()

		router := gin.Default()
		rec := httptest.NewRecorder()

		userHandler := delivery.NewUserHandler(router, userUseCase)
		router.POST("/user/register", userHandler.Register)

		// do
		reqBody, err := json.Marshal(tempMockRegisterUser)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/user/register", strings.NewReader(string(reqBody)))
		router.ServeHTTP(rec, req)

		var res domain.RegisteredUser
		err = json.Unmarshal(rec.Body.Bytes(), &res)
		assert.NoError(t, err)

		_, err = govalidator.ValidateStruct(tempMockRegisterUser)

		// assert
		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "age: 5 does not validate as range(8|63)", res.Message)
	})
}

func TestLoginUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUserUseCase := new(mocksUseCase.UserUseCase)
	userUseCase := userUseCase.NewUserUseCase(mockUserUseCase)

	// t.Run("should success login user", func(t *testing.T) {
	// 	// prepare
	// 	tempMockLoginUser := domain.User{
	// 		Email:    "johndoe@example.com",
	// 		Password: "secret",
	// 	}
	// 	expected := domain.LoggedInUser{
	// 		Status:  "success",
	// 		Message: "user registration has been successful",
	// 		Data: domain.Token{
	// 			Token: "your-token",
	// 		},
	// 	}

	// 	mockUserUseCase.On("Login", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil).Once()

	// 	router := gin.Default()
	// 	rec := httptest.NewRecorder()

	// 	userHandler := delivery.NewUserHandler(router, userUseCase)
	// 	router.POST("/user/login", userHandler.Login)

	// 	// do
	// 	reqBody, err := json.Marshal(tempMockLoginUser)
	// 	assert.NoError(t, err)
	// 	req := httptest.NewRequest(http.MethodPost, "/user/login", strings.NewReader(string(reqBody)))
	// 	router.ServeHTTP(rec, req)

	// 	var res domain.LoggedInUser
	// 	err = json.Unmarshal(rec.Body.Bytes(), &res)
	// 	assert.NoError(t, err)

	// 	// assert
	// 	assert.Equal(t, http.StatusOK, rec.Code)
	// 	assert.Equal(t, expected.Message, res.Message)
	// })

	t.Run("should fail login user with invalid email", func(t *testing.T) {
		// prepare
		tempMockLoginUser := domain.User{
			Email:    "jhondoe@example.com",
			Password: "secret",
		}

		mockUserUseCase.On("Login", mock.Anything, mock.AnythingOfType("*domain.User")).Return(errors.New("the email you entered are not registered")).Once()

		router := gin.Default()
		rec := httptest.NewRecorder()

		userHandler := delivery.NewUserHandler(router, userUseCase)
		router.POST("/user/login", userHandler.Login)

		// do
		reqBody, err := json.Marshal(tempMockLoginUser)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/user/login", strings.NewReader(string(reqBody)))
		router.ServeHTTP(rec, req)

		var res domain.RegisteredUser
		err = json.Unmarshal(rec.Body.Bytes(), &res)
		assert.NoError(t, err)

		_, err = govalidator.ValidateStruct(tempMockLoginUser)

		// assert
		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "the email you entered are not registered", res.Message)
	})

	t.Run("should fail login user with invalid password", func(t *testing.T) {
		// prepare
		tempMockLoginUser := domain.User{
			Email:    "jhondoe@example.com",
			Password: "secret",
		}

		mockUserUseCase.On("Login", mock.Anything, mock.AnythingOfType("*domain.User")).Return(errors.New("the password you entered are wrong")).Once()

		router := gin.Default()
		rec := httptest.NewRecorder()

		userHandler := delivery.NewUserHandler(router, userUseCase)
		router.POST("/user/login", userHandler.Login)

		// do
		reqBody, err := json.Marshal(tempMockLoginUser)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/user/login", strings.NewReader(string(reqBody)))
		router.ServeHTTP(rec, req)

		var res domain.RegisteredUser
		err = json.Unmarshal(rec.Body.Bytes(), &res)
		assert.NoError(t, err)

		_, err = govalidator.ValidateStruct(tempMockLoginUser)

		// assert
		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "the password you entered are wrong", res.Message)
	})
}
