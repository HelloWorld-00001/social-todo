package cmd

import (
	"fmt"
	goService "github.com/200Lab-Education/go-sdk"
	"github.com/200Lab-Education/go-sdk/plugin/storage/sdkgorm"
	"github.com/coderconquerer/go-login-app/configs"
	"github.com/coderconquerer/go-login-app/docs"
	accBuc "github.com/coderconquerer/go-login-app/internal/account/BusinessUseCases"
	accountHdl "github.com/coderconquerer/go-login-app/internal/account/Handler"
	accStorage "github.com/coderconquerer/go-login-app/internal/account/Storage"
	"github.com/coderconquerer/go-login-app/internal/common"
	"github.com/coderconquerer/go-login-app/internal/components/tokenProviders/jwtProvider"
	"github.com/coderconquerer/go-login-app/internal/components/uploadProvider"
	"github.com/coderconquerer/go-login-app/internal/file/BusinessUseCases"
	"github.com/coderconquerer/go-login-app/internal/file/Handler"
	"github.com/coderconquerer/go-login-app/internal/middleware"
	todoBuc "github.com/coderconquerer/go-login-app/internal/todoItem/BusinessUseCases"
	todoHdl "github.com/coderconquerer/go-login-app/internal/todoItem/Handler"
	todoStorage "github.com/coderconquerer/go-login-app/internal/todoItem/Storage"
	userBuc "github.com/coderconquerer/go-login-app/internal/user/BusinessUseCases"
	userHdl "github.com/coderconquerer/go-login-app/internal/user/Handler"
	userStorage "github.com/coderconquerer/go-login-app/internal/user/Storage"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
	"os"
)

func NewServices() goService.Service {
	service := goService.New(
		goService.WithName("social-todo-app"),
		goService.WithVersion("1.0.0"),
		goService.WithInitRunnable(sdkgorm.NewGormDB("main-db", common.DbMainName)),
	)

	return service
}

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Start social todo service",
	Run: func(cmd *cobra.Command, args []string) {
		service := NewServices()
		serviceLog := service.Logger("service-todo")

		if err := service.Init(); err != nil {
			serviceLog.Fatalln(err)
		}

		service.HTTPServer().AddHandler(func(engine *gin.Engine) {
			cfg := configs.Load()
			awsCfg := configs.LoadAWSConfig()
			database := service.MustGet(common.DbMainName).(*gorm.DB)

			// init services
			tokenProvider := jwtProvider.GetNewJwtProvider(cfg.JwtConfig.Prefix, cfg.JwtConfig.SecretKey)
			accStore := accStorage.GetNewMySQLConnection(database)
			registerBsn := accBuc.GetNewRegisterAccountLogic(accStore)
			loginBsn := accBuc.GetNewLoginLogic(accStore, tokenProvider, 60*60*24*30)
			disableBsn := accBuc.GetNewDisableAccountLogic(accStore)
			accountHandler := accountHdl.NewAccountHandler(loginBsn, registerBsn, disableBsn)

			// to do service
			todoStore := todoStorage.GetNewMySQLConnection(database)
			getTodoDetailBz := todoBuc.GetNewGetTodoDetailLogic(todoStore)
			getTodoListBz := todoBuc.GetNewGetTodoListLogic(todoStore)
			createTodoListBz := todoBuc.GetNewCreateTodoLogic(todoStore)
			deleteTodoListBz := todoBuc.GetNewDeleteTodoItemLogic(todoStore)
			todoHandler := todoHdl.NewTodoHandler(getTodoDetailBz, getTodoListBz, createTodoListBz, deleteTodoListBz)

			// user service
			userStore := userStorage.GetNewMySQLConnection(database)
			getUserProfileBz := userBuc.GetNewFindUserLogic(userStore)
			userHandler := userHdl.NewUserHandler(getUserProfileBz)

			// aws services
			s3Provider := uploadProvider.NewS3ProviderWithConfig(awsCfg)
			uploadBsn := BusinessUseCases.GetNewUploadFileLogic(todoStore, userStore, s3Provider)
			uploadHandler := Handler.NewUploadHandler(uploadBsn)
			swaggerSetup()

			// init middleware
			authAdmin := middleware.RequireAuth(tokenProvider, accStore, common.AdminRole.ToString())
			authUser := middleware.RequireAuth(tokenProvider, accStore, []string{common.AdminRole.ToString(), common.UserRole.ToString()}...) // multi-role access

			//engine.Use(gin.Logger())
			//engine.Use(middleware.CustomRecovery()) // custom middleware
			//auth := engine.Group("/")b
			store := cookie.NewStore([]byte("super-secret-key"))

			v1 := engine.Group("/v1/api")
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
				v1.POST("/upload", authUser, uploadHandler.UploadImage())

			}
			engine.Use(sessions.Sessions("mysession", store))
			engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
			engine.GET("/v1/api/todo/ping", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "pong",
				})
			})
		})

		if err := service.Start(); err != nil {
			serviceLog.Fatalln(err)
		}

	},
}

func Execute() {
	rootCmd.AddCommand(outEnvCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
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
