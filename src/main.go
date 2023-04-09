package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gusrylmubarok/mygram-backend/src/config"
	"github.com/joho/godotenv"

	docs "github.com/gusrylmubarok/mygram-backend/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	commentDelivery "github.com/gusrylmubarok/mygram-backend/src/modules/comment/delivery/http"
	commentRepository "github.com/gusrylmubarok/mygram-backend/src/modules/comment/repository/postgres"
	commentUseCase "github.com/gusrylmubarok/mygram-backend/src/modules/comment/usecase"
	photoDelivery "github.com/gusrylmubarok/mygram-backend/src/modules/photo/delivery/http"
	photoRepository "github.com/gusrylmubarok/mygram-backend/src/modules/photo/repository/postgres"
	photoUseCase "github.com/gusrylmubarok/mygram-backend/src/modules/photo/usecase"
	socialMediaDelivery "github.com/gusrylmubarok/mygram-backend/src/modules/socialmedia/delivery/http"
	socialMediaRepository "github.com/gusrylmubarok/mygram-backend/src/modules/socialmedia/repository/postgres"
	socialMediaUseCase "github.com/gusrylmubarok/mygram-backend/src/modules/socialmedia/usecase"
	userDelivery "github.com/gusrylmubarok/mygram-backend/src/modules/user/delivery/http"
	userRepository "github.com/gusrylmubarok/mygram-backend/src/modules/user/repository/postgres"
	userUseCase "github.com/gusrylmubarok/mygram-backend/src/modules/user/usecase"
)

// @title			mygram backend
// @version 		1.0.0
// @description 	mygram is a free photo sharing app written in Go. People can share, view, and comment photos by everyone. Anyone can create an account by registering an email address and creating a username.
// @termOfService 	http://swagger.io/terms/
// @contact.name 	gusrylmubarok
// @contact.email 	gusrylmubarok@gmail.com
// @license.name 	MIT License
// @license.url 	https://opensource.org/licenses/MIT
// @host 			localhost:8080
// @BasePath 		/

// @securityDefinitions.apikey  Bearer
// @in                          header
// @name                        Authorization
// @description					Description for what is this security definition being used

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	db := config.ConnectDB()

	routers := gin.Default()
	routers.Use(func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Content-Type", "application/json")
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Max-Age", "86400")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, UPDATE")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(200)
		} else {
			ctx.Next()
		}
	})
	routers.Static("/public", "./public")

	userRepository := userRepository.NewUserRepository(db)
	userUseCase := userUseCase.NewUserUseCase(userRepository)

	userDelivery.NewUserHandler(routers, userUseCase)

	photoRepository := photoRepository.NewPhotoRepository(db)
	photoUseCase := photoUseCase.NewPhotoUseCase(photoRepository)

	photoDelivery.NewPhotoHandler(routers, photoUseCase)

	commentRepository := commentRepository.NewCommentRepository(db)
	commentUseCase := commentUseCase.NewCommentUseCase(commentRepository)

	commentDelivery.NewCommentHandler(routers, commentUseCase, photoUseCase)

	socialMediaRepository := socialMediaRepository.NewSocialMediaRepository(db)
	socialMediaUseCase := socialMediaUseCase.NewSocialMediaUseCase(socialMediaRepository)

	socialMediaDelivery.NewSocialMediaHandler(routers, socialMediaUseCase)

	docs.SwaggerInfo.BasePath = "/api/v1"
	routers.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	port := os.Getenv("PORT")

	if len(os.Args) > 1 {
		reqPort := os.Args[1]

		if reqPort != "" {
			port = reqPort
		}
	}

	if port == "" {
		port = "8080"
	}

	routers.Run(":" + port)
}
