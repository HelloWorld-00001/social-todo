package servicesetup

import (
	goService "github.com/200Lab-Education/go-sdk"
	"github.com/coderconquerer/social-todo"
	"github.com/coderconquerer/social-todo/docs"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func StartHttpServer(service goService.Service) {
	service.HTTPServer().AddHandler(func(engine *gin.Engine) {
		// session middleware
		store := cookie.NewStore([]byte("super-secret-key"))
		engine.Use(sessions.Sessions("mysession", store))

		// app routes
		social_todo.SetupRoutes(engine.Group(""), service)

		// swagger docs
		swaggerSetup()
		engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		// health check
		engine.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})
	})
}

func swaggerSetup() {
	docs.SwaggerInfo.Title = "Login API"
	docs.SwaggerInfo.Description = "API documentation for login system"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}
}
