package servicesetup

import (
	goService "github.com/200Lab-Education/go-sdk"
	"github.com/200Lab-Education/go-sdk/plugin/storage/sdkgorm"
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/configs"
	"github.com/coderconquerer/social-todo/plugin/redis"
	"github.com/coderconquerer/social-todo/plugin/rpc_caller"
	"github.com/coderconquerer/social-todo/plugin/tokenprovider/jwtProvider"
	"github.com/coderconquerer/social-todo/plugin/uploadprovider/s3provider"
	"github.com/coderconquerer/social-todo/pubsub"
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
		goService.WithInitRunnable(jwtprovider.GetNewJwtProvider(cfg.JwtConfig.JwtPrefix, cfg.JwtConfig.SecretKey)),
		goService.WithInitRunnable(s3provider.NewS3ProviderWithConfig(awsCfg)),
		goService.WithInitRunnable(pubsub.NewLocalPubsub(common.PluginPubSub)),
		goService.WithInitRunnable(rpc_caller.NewRpcCaller(common.PluginRPC)),
		goService.WithInitRunnable(redis.NewRedisDB("redisDb", common.PluginRedis)),
	)

	return service
}
