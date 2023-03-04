package routes

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/forkyid/go-utils/v1/middleware"
	"github.com/gin-gonic/gin"
	"go-rest-api/docs"
	"go-rest-api/src/connection"
	"gorm.io/gorm"

	authController "go-rest-api/src/controller/v1/auth"
	accountController "go-rest-api/src/controller/v1/account"
	locationController "go-rest-api/src/controller/v1/location"

	accountRepository "go-rest-api/src/repository/v1/account"
	locationRepository "go-rest-api/src/repository/v1/location"

	accountService "go-rest-api/src/service/v1/account"
	locationService "go-rest-api/src/service/v1/location"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var master *gorm.DB
var router = gin.Default()

type DB struct {
	Master *gorm.DB
}

func Run() {	
	godotenv.Load()
	RouterSetup()
	router.Run(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")))
}

func RouterSetup() *gin.Engine {
	// set up
	router.SetTrustedProxies(nil)
	middleware := middleware.Middleware{}
	router.Use(middleware.CORS)

	// swagger
	docs.SwaggerInfo.Title = "Phincon Attendance Rest API"
	docs.SwaggerInfo.Description = "Phincon Attendance Rest API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = os.Getenv("SWAGGER_HOST")
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// database connection (type *gorm.DB)
	master = connection.DBMaster()

	// repository
	accountRepo := accountRepository.NewRepository(connection.DB{
		Master: master,
	})
	locationRepo := locationRepository.NewRepository(connection.DB{
		Master: master,
	})

	// service
	accountSvc := accountService.NewService(accountRepo)
	locationSvc := locationService.NewService(locationRepo)
	
	// controller
	authController := authController.NewController(accountSvc)
	accountController := accountController.NewController(accountSvc)
	locationController := locationController.NewController(locationSvc)

	// endpoint
	v1 := router.Group("v1")

	auth := v1.Group("auth")
	auth.POST("", authController.Login)
	auth.PATCH("forgot", authController.ForgotPassword)

	accounts := v1.Group("accounts")
	accounts.GET("", accountController.Get)
	accounts.POST("register", accountController.Register)
	accounts.PATCH("", accountController.Update)
	accounts.DELETE("", accountController.Delete)

	location := v1.Group("locations")
	location.GET("", locationController.Get)
	location.POST("", locationController.Create)
	location.PATCH("", locationController.Update)
	location.DELETE("", locationController.Delete)

	return router
}
