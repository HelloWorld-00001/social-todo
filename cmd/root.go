package cmd

import (
	"fmt"
	goService "github.com/200Lab-Education/go-sdk"
	"github.com/200Lab-Education/go-sdk/plugin/storage/sdkgorm"
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/configs"
	"github.com/coderconquerer/social-todo/docs"
	"github.com/coderconquerer/social-todo/middleware"
	accBuc "github.com/coderconquerer/social-todo/module/account/BusinessUseCases"
	accountHdl "github.com/coderconquerer/social-todo/module/account/Handler"
	accStorage "github.com/coderconquerer/social-todo/module/account/Storage"
	"github.com/coderconquerer/social-todo/module/file/BusinessUseCases"
	"github.com/coderconquerer/social-todo/module/file/Handler"
	todoBuc "github.com/coderconquerer/social-todo/module/todoItem/BusinessUseCases"
	todoHdl "github.com/coderconquerer/social-todo/module/todoItem/Handler"
	"github.com/coderconquerer/social-todo/module/todoItem/Repository"
	todoStorage "github.com/coderconquerer/social-todo/module/todoItem/Storage"
	userBuc "github.com/coderconquerer/social-todo/module/user/BusinessUseCases"
	userHdl "github.com/coderconquerer/social-todo/module/user/Handler"
	userStorage "github.com/coderconquerer/social-todo/module/user/Storage"
	BusinessUseCases2 "github.com/coderconquerer/social-todo/module/userReactItem/BusinessUseCases"
	Handler2 "github.com/coderconquerer/social-todo/module/userReactItem/Handler"
	"github.com/coderconquerer/social-todo/module/userReactItem/Storage"
	tokenPlugin "github.com/coderconquerer/social-todo/plugin/tokenProviders"
	"github.com/coderconquerer/social-todo/plugin/tokenProviders/jwtProvider"
	uploadPlugin "github.com/coderconquerer/social-todo/plugin/uploadProvider"
	"github.com/coderconquerer/social-todo/plugin/uploadProvider/s3provider"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
	"os"
)

var (
	cfg    = configs.Load()
	awsCfg = configs.LoadAWSConfig()
)

func NewServices() goService.Service {
	service := goService.New(
		goService.WithName("social-todo-app"),
		goService.WithVersion("1.0.0"),
		goService.WithInitRunnable(sdkgorm.NewGormDB("main-db", common.DbMainName)),
		goService.WithInitRunnable(jwtProvider.GetNewJwtProvider(cfg.JwtConfig.JwtPrefix, cfg.JwtConfig.SecretKey)),
		goService.WithInitRunnable(s3provider.NewS3ProviderWithConfig(awsCfg)),
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
			database := service.MustGet(common.DbMainName).(*gorm.DB)
			tokenProvider := service.MustGet(cfg.JwtConfig.JwtPrefix).(tokenPlugin.TokenProvider)
			s3Provider := service.MustGet(awsCfg.S3Prefix).(uploadPlugin.UploadProvider)

			// init services
			accStore := accStorage.GetNewMySQLConnection(database)
			registerBsn := accBuc.GetNewRegisterAccountLogic(accStore)
			loginBsn := accBuc.GetNewLoginLogic(accStore, tokenProvider, 60*60*24*30)
			disableBsn := accBuc.GetNewDisableAccountLogic(accStore)
			accountHandler := accountHdl.NewAccountHandler(loginBsn, registerBsn, disableBsn)

			// reaction service
			reactionStore := Storage.GetNewMySQLConnection(database)
			reactBz := BusinessUseCases2.GetNewReactTodoItemLogic(reactionStore)
			unReactBz := BusinessUseCases2.GetNewUnreactTodoItemLogic(reactionStore)
			listRUBz := BusinessUseCases2.GetNewGetListReactedUsersLogic(reactionStore)
			reactHandler := Handler2.NewReactionHandler(reactBz, unReactBz, listRUBz)

			// to do service
			todoStore := todoStorage.GetNewMySQLConnection(database)
			getTodoDetailBz := todoBuc.GetNewGetTodoDetailLogic(todoStore)
			getTodoListRepo := Repository.GetNewTodoListWithReactRepo(todoStore, reactionStore)
			getTodoListBz := todoBuc.GetNewGetTodoListLogic(getTodoListRepo)
			createTodoListBz := todoBuc.GetNewCreateTodoLogic(todoStore)
			deleteTodoListBz := todoBuc.GetNewDeleteTodoItemLogic(todoStore)
			todoHandler := todoHdl.NewTodoHandler(getTodoDetailBz, getTodoListBz, createTodoListBz, deleteTodoListBz)

			// user service
			userStore := userStorage.GetNewMySQLConnection(database)
			getUserProfileBz := userBuc.GetNewFindUserLogic(userStore)
			userHandler := userHdl.NewUserHandler(getUserProfileBz)

			// aws services
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

				react := v1.Group("/react")
				{
					react.GET("/:todo_id", authUser, reactHandler.GetListReactedUsers())
					react.POST("", authUser, reactHandler.ReactItem())
					react.DELETE("", authUser, reactHandler.UnreactTodoItem())
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
