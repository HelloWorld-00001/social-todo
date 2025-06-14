package main

import (
	"github.com/coderconquerer/go-login-app/docs"
	_ "github.com/coderconquerer/go-login-app/docs"
	al "github.com/coderconquerer/go-login-app/internal/accountLogic"
	"github.com/coderconquerer/go-login-app/pkg/config"
	db "github.com/coderconquerer/go-login-app/pkg/database"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.Load()
	database, err := db.Connect(cfg)
	if err != nil {
		panic(err)
	}

	repo := al.NewRepository(database)
	service := al.NewService(repo)
	controller := al.NewController(service)
	swaggerSetup()

	r := gin.Default()

	//auth := r.Group("/")b
	store := cookie.NewStore([]byte("super-secret-key"))
	r.Use(sessions.Sessions("mysession", store))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/register", controller.RegisterUser)
	r.POST("/login", controller.LoginUser)
	r.GET("/api_count", controller.GetAPICount)
	err = r.Run(":8080")
	if err != nil {
		return
	}

}

func swaggerSetup() {
	docs.SwaggerInfo.Title = "Login API"
	docs.SwaggerInfo.Description = "API documentation for login system"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}
}
