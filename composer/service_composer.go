package composer

import (
	serviceCtx "github.com/200Lab-Education/go-sdk"
	"github.com/coderconquerer/social-todo/cmd/registerservice"
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/configs"
	"github.com/coderconquerer/social-todo/grpc/contract"
	authBusiness "github.com/coderconquerer/social-todo/module/authentication/business"
	authStorage "github.com/coderconquerer/social-todo/module/authentication/storage"
	"github.com/coderconquerer/social-todo/module/authentication/transport"
	uploadBusiness "github.com/coderconquerer/social-todo/module/file/business"
	"github.com/coderconquerer/social-todo/module/file/transport"
	"github.com/coderconquerer/social-todo/module/todo/business"
	"github.com/coderconquerer/social-todo/module/todo/repository"
	restapi "github.com/coderconquerer/social-todo/module/todo/storage/grpc"
	"github.com/coderconquerer/social-todo/module/todo/storage/mysql"
	"github.com/coderconquerer/social-todo/module/todo/transport/api"
	reactionBusiness "github.com/coderconquerer/social-todo/module/todotaskreaction/business"
	reactionStorage "github.com/coderconquerer/social-todo/module/todotaskreaction/storage"
	reactionAPI "github.com/coderconquerer/social-todo/module/todotaskreaction/transport/api"
	userBusiness "github.com/coderconquerer/social-todo/module/user/business"
	userStorage "github.com/coderconquerer/social-todo/module/user/storage"
	"github.com/coderconquerer/social-todo/module/user/transport"
	tokenPlugin "github.com/coderconquerer/social-todo/plugin/tokenProviders"
	uploadPlugin "github.com/coderconquerer/social-todo/plugin/uploadProvider"
	"github.com/coderconquerer/social-todo/pubsub"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
	"log"
)

var (
	cfg    = configs.Load()
	awsCfg = configs.LoadAWSConfig()
)

func GetTodoAPIService(sCtx serviceCtx.ServiceContext) api.TodoAPI {

	database := sCtx.MustGet(common.DbMainName).(*gorm.DB)

	storage := mysql.GetNewMySQLConnection(database)
	conn, err := grpc.NewClient(
		"0.0.0.0:8082",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatalln("Failed to connect to grpc server:", err)
	}
	//	TODO implement gateway
	// Create client stub
	client := contract.NewItemReactServiceClient(conn)
	rpcClient := restapi.NewRpcClient(client)
	repo := repository.GetNewTodoListWithReactRepo(storage, rpcClient)
	bz := business.NewTodoBusiness(storage, repo)

	return api.NewTodoAPI(bz)
}

func GetAuthenticationAPIService(sCtx serviceCtx.ServiceContext) authAPI.AuthenticationAPI {

	database := sCtx.MustGet(common.DbMainName).(*gorm.DB)
	tokenProvider := sCtx.MustGet(cfg.JwtConfig.JwtPrefix).(tokenPlugin.TokenProvider)

	storage := authStorage.GetNewMySQLConnection(database)
	// init services
	bz := authBusiness.NewAuthenticationBusiness(storage, tokenProvider, 60*60*24*30)
	// todo move expire time to config

	return authAPI.NewAuthenticationAPI(bz)
}

func GetUserAPIService(sCtx serviceCtx.ServiceContext) userApi.UserAPI {

	database := sCtx.MustGet(common.DbMainName).(*gorm.DB)

	storage := userStorage.GetNewMySQLConnection(database)
	bz := userBusiness.NewUserBusiness(storage)
	return userApi.NewUserAPI(bz)
}

func GetTodoReactionService(sCtx serviceCtx.ServiceContext) reactionBusiness.ReactionBusiness {

	database := sCtx.MustGet(common.DbMainName).(*gorm.DB)
	ps := sCtx.MustGet(common.PluginPubSub).(pubsub.PubSub)

	storage := reactionStorage.GetNewMySQLConnection(database)

	rabbitMQ := registerservice.NewRabbitMQPublisher("amqp://guest:guest@localhost:5672/", "RabbitMQ_Test")
	return reactionBusiness.NewReactionBusiness(storage, ps, rabbitMQ)
}

func GetTodoReactionAPIService(sCtx serviceCtx.ServiceContext) reactionAPI.ReactionTodoAPI {

	database := sCtx.MustGet(common.DbMainName).(*gorm.DB)
	ps := sCtx.MustGet(common.PluginPubSub).(pubsub.PubSub)

	storage := reactionStorage.GetNewMySQLConnection(database)

	rabbitMQ := registerservice.NewRabbitMQPublisher("amqp://guest:guest@localhost:5672/", "RabbitMQ_Test")
	bz := reactionBusiness.NewReactionBusiness(storage, ps, rabbitMQ)

	return reactionAPI.NewReactionTodoAPI(bz)
}

func GetUploadFileAPIService(sCtx serviceCtx.ServiceContext) uploadImageAPI.UploadImageAPI {
	database := sCtx.MustGet(common.DbMainName).(*gorm.DB)
	s3Provider := sCtx.MustGet(awsCfg.S3Prefix).(uploadPlugin.UploadProvider)

	todoStorage := mysql.GetNewMySQLConnection(database)
	userStr := userStorage.GetNewMySQLConnection(database)
	bz := uploadBusiness.NewUploadImageLogic(todoStorage, userStr, s3Provider)
	return uploadImageAPI.NewUploadImageAPI(bz)
}
