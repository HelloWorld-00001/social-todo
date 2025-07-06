package main

import (
	"fmt"
	"github.com/coderconquerer/go-login-app/docs"
	_ "github.com/coderconquerer/go-login-app/docs"
	TodoHandler "github.com/coderconquerer/go-login-app/internal/TodoItem/Handler"
	accBuc "github.com/coderconquerer/go-login-app/internal/account/BusinessUseCases"
	AccountHandler "github.com/coderconquerer/go-login-app/internal/account/Handler"
	accStorage "github.com/coderconquerer/go-login-app/internal/account/Storage"
	"github.com/coderconquerer/go-login-app/internal/common"
	todoBuc "github.com/coderconquerer/go-login-app/internal/todoItem/BusinessUseCases"
	todoStorage "github.com/coderconquerer/go-login-app/internal/todoItem/Storage"
	userBuc "github.com/coderconquerer/go-login-app/internal/user/BusinessUseCases"
	userHandler "github.com/coderconquerer/go-login-app/internal/user/Handler"
	userStorage "github.com/coderconquerer/go-login-app/internal/user/Storage"

	"github.com/coderconquerer/go-login-app/internal/components/tokenProviders/jwtProvider"
	"github.com/coderconquerer/go-login-app/internal/middleware"
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
	database, err := todoStorage.GetMySQLConnection(cfg.DbConfig)
	if err != nil {
		fmt.Println("Cannot connect to database")
		fmt.Println(cfg)
		panic(err)
	}

	// init services
	tokenProvider := jwtProvider.GetNewJwtProvider(cfg.JwtConfig.Prefix, cfg.JwtConfig.SecretKey)
	accStore := accStorage.GetNewMySQLConnection(database)
	registerBsn := accBuc.GetNewRegisterAccountLogic(accStore)
	loginBsn := accBuc.GetNewLoginLogic(accStore, tokenProvider, 60*60*24*30)
	disableBsn := accBuc.GetNewDisableAccountLogic(accStore)
	accountHandler := AccountHandler.NewAccountHandler(loginBsn, registerBsn, disableBsn)

	// to do service
	todoStore := todoStorage.GetNewMySQLConnection(database)
	getTodoDetailBz := todoBuc.GetNewGetTodoDetailLogic(todoStore)
	getTodoListBz := todoBuc.GetNewGetTodoListLogic(todoStore)
	createTodoListBz := todoBuc.GetNewCreateTodoLogic(todoStore)
	deleteTodoListBz := todoBuc.GetNewDeleteTodoItemLogic(todoStore)
	todoHandler := TodoHandler.NewTodoHandler(getTodoDetailBz, getTodoListBz, createTodoListBz, deleteTodoListBz)

	// user service
	userStore := userStorage.GetNewMySQLConnection(database)
	getUserProfileBz := userBuc.GetNewFindUserLogic(userStore)
	userHandler := userHandler.NewUserHandler(getUserProfileBz)

	swaggerSetup()

	r := gin.Default()

	// init middleware
	authAdmin := middleware.RequireAuth(tokenProvider, accStore, common.AdminRole.ToString())
	authUser := middleware.RequireAuth(tokenProvider, accStore, []string{common.AdminRole.ToString(), common.UserRole.ToString()}...) // multi-role access

	//r.Use(gin.Logger())
	//r.Use(middleware.CustomRecovery()) // custom middleware
	//auth := r.Group("/")b
	store := cookie.NewStore([]byte("super-secret-key"))

	v1 := r.Group("/v1/api")
	{
		todoRoutes := v1.Group("/todo")
		{
			todoRoutes.GET("", todoHandler.GetToDoList())
			todoRoutes.GET("/:id", todoHandler.GetTodoDetail())
			//todoRoutes.PUT("/:id", Handler.GetTodoDetail(database))
			todoRoutes.DELETE("/:id", authUser, todoHandler.DeleteTodoItem())
			todoRoutes.POST("/", authUser, todoHandler.CreateTodoItem())
		}
		v1.GET("/profile", authUser, userHandler.GetUserProfile())
		v1.POST("/login", accountHandler.Login())
		v1.POST("/register", accountHandler.RegisterAccount())
		v1.POST("/disable", authAdmin, accountHandler.DisableAccount())

	}
	r.Use(sessions.Sessions("mysession", store))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
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
