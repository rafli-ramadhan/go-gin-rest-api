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
	accountRepository "go-rest-api/src/repository/v1/account"
	accountService "go-rest-api/src/service/v1/account"

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
	docs.SwaggerInfo.Title = "Go Rest API"
	docs.SwaggerInfo.Description = "Go Rest API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = os.Getenv("SWAGGER_HOST")
	docs.SwaggerInfo.Schemes = []string{"https", "http"}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// database connection (type *gorm.DB)
	master = connection.DBMaster()

	// repository
	accountRepository := accountRepository.NewRepository(connection.DB{
		Master: master,
	})

	// service
	accountService := accountService.NewService(accountRepository)
	
	// controller
	authController := authController.NewController(accountService)
	accountController := accountController.NewController(accountService)

	// endpoint
	v1 := router.Group("v1")

	auth := v1.Group("auth")
	auth.POST("", authController.Login)
	auth.GET("self", authController.AuthSelf)

	accounts := v1.Group("accounts")
	accounts.POST("register", accountController.Register)
	accounts.PATCH("", accountController.Update)
	accounts.DELETE("", accountController.Delete)

	return router
}
