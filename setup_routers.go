package social_todo

import (
	serviceCtx "github.com/200Lab-Education/go-sdk"
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/composer"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(route *gin.RouterGroup, sCtx serviceCtx.ServiceContext) {

	authUser := composer.GetAuthenticationMiddleWare(sCtx, common.UserRole)
	authAdmin := composer.GetAuthenticationMiddleWare(sCtx, common.AdminRole)

	todoAPI := composer.GetTodoAPIService(sCtx)
	authAPI := composer.GetAuthenticationAPIService(sCtx)
	reactAPI := composer.GetTodoReactionAPIService(sCtx)
	userAPI := composer.GetUserAPIService(sCtx)
	uploadImageAPI := composer.GetUploadFileAPIService(sCtx)
	authGrpcAPI := composer.GetAuthenticationGrpcAPIService(sCtx)

	v1 := route.Group("/v1/api")
	{
		todoRoutes := v1.Group("/todo")
		{
			todoRoutes.GET("", authUser, todoAPI.GetToDoList())
			todoRoutes.GET("/:id", todoAPI.GetTodoDetail())
			todoRoutes.DELETE("/:id", authUser, todoAPI.DeleteTodoItem())
			todoRoutes.POST("/", authUser, todoAPI.CreateTodoItem())
		}
		react := v1.Group("/react")
		{
			react.GET("/:todo_id", authUser, reactAPI.GetListReactedUsers())
			react.POST("", authUser, reactAPI.ReactItem())
			react.DELETE("", authUser, reactAPI.UnreactTodoItem())
		}

		auth := v1.Group("/auth")
		{
			auth.POST("/login", authAPI.Login())
			auth.POST("/register", authAPI.RegisterAccount())
		}

		rpcAuth := v1.Group("/rpc/auth")
		{
			rpcAuth.POST("/login", authGrpcAPI.Login())
			rpcAuth.POST("/register", authGrpcAPI.RegisterAccount())
		}

		v1.GET("/profile", authUser, userAPI.GetUserProfile())
		v1.POST("/disable", authAdmin, authAPI.DisableAccount())
		v1.POST("/upload", authUser, uploadImageAPI.UploadImage())

	}
}
