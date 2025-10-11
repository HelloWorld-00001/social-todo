package composer

import (
	serviceCtx "github.com/200Lab-Education/go-sdk"
	"github.com/coderconquerer/social-todo/cache"
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/configs"
	"github.com/coderconquerer/social-todo/middleware"
	accStorage "github.com/coderconquerer/social-todo/module/authentication/storage"
	tokenPlugin "github.com/coderconquerer/social-todo/plugin/tokenprovider"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAuthenticationMiddleWare(sCtx serviceCtx.ServiceContext, roleLevel common.Role) func(c *gin.Context) {
	database := sCtx.MustGet(common.DbMainName).(*gorm.DB)
	cfg := configs.Load()
	tokenProvider := sCtx.MustGet(cfg.JwtConfig.JwtPrefix).(tokenPlugin.TokenProvider)

	// redis cache
	accStore := accStorage.GetNewMySQLConnection(database)
	redisService := cache.NewRedisCache(sCtx)
	userCache := cache.NewUserCaching(redisService, accStore)
	// init middleware
	switch roleLevel {
	case common.AdminRole:
		return middleware.RequireAuth(tokenProvider, userCache, common.AdminRole.ToString())
	case common.UserRole:
		return middleware.RequireAuth(tokenProvider, userCache, []string{common.AdminRole.ToString(), common.UserRole.ToString()}...) // multi-role access
	default:
		return nil
	}
}
