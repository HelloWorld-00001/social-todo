package main

import (
	"fmt"
	"github.com/coderconquerer/go-login-app/docs"
	_ "github.com/coderconquerer/go-login-app/docs"
	db "github.com/coderconquerer/go-login-app/internal/Storage"
	"github.com/coderconquerer/go-login-app/internal/TodoHandler"
	al "github.com/coderconquerer/go-login-app/internal/accountLogic"
	"github.com/coderconquerer/go-login-app/pkg/config"
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
	database, err := db.GetMySQLConnection(cfg)
	if err != nil {
		fmt.Println("Cannot connect to database")
		fmt.Println(cfg)
		panic(err)
	}

	repo := al.NewRepository(nil)
	service := al.NewService(repo)
	controller := al.NewController(service)
	swaggerSetup()

	r := gin.Default()

	//auth := r.Group("/")b
	store := cookie.NewStore([]byte("super-secret-key"))

	v1 := r.Group("/v1/api")
	{
		todoRoutes := v1.Group("/todo")
		{
			todoRoutes.GET("", TodoHandler.GetToDoList(database))
			todoRoutes.GET("/:id", TodoHandler.GetTodoDetail(database))
			//todoRoutes.PUT("/:id", TodoHandler.GetTodoDetail(database))
			todoRoutes.DELETE("/:id", TodoHandler.DeleteTodoItem(database))
			todoRoutes.POST("/", TodoHandler.CreateTodoItem(database))
		}

	}
	r.Use(sessions.Sessions("mysession", store))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/register", controller.RegisterUser)
	r.POST("/login", controller.LoginUser)
	r.GET("/api_count", controller.GetAPICount)
	r.GET("/v1/api/todo/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
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
