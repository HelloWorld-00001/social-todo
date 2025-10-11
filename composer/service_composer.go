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
	authGprcBusiness "github.com/coderconquerer/social-todo/module/authenticationrpc/business"
	"github.com/coderconquerer/social-todo/module/authenticationrpc/storage"
	"github.com/coderconquerer/social-todo/module/authenticationrpc/transport"
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
	tokenPlugin "github.com/coderconquerer/social-todo/plugin/tokenprovider"
	uploadPlugin "github.com/coderconquerer/social-todo/plugin/uploadprovider"
	"github.com/coderconquerer/social-todo/pubsub"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
	"log"
)

var (
	cfg        = configs.Load()
	awsCfg     = configs.LoadAWSConfig()
	grpcConfig = configs.LoadGrpcPort()
)

func GetTodoAPIService(sCtx serviceCtx.ServiceContext) api.TodoAPI {

	database := sCtx.MustGet(common.DbMainName).(*gorm.DB)

	strg := mysql.GetNewMySQLConnection(database)
	conn, err := grpc.NewClient(
		"localhost:"+grpcConfig.TodoReactionPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatalln("Failed to connect to grpc server:", err)
	}
	//	TODO implement gateway
	// Create client stub
	client := contract.NewItemReactServiceClient(conn)
	rpcClient := restapi.NewRpcClient(client)
	repo := repository.GetNewTodoListWithReactRepo(strg, rpcClient)
	bz := business.NewTodoBusiness(strg, repo)

	return api.NewTodoAPI(bz)
}

func GetAuthenticationAPIService(sCtx serviceCtx.ServiceContext) authAPI.AuthenticationAPI {

	database := sCtx.MustGet(common.DbMainName).(*gorm.DB)
	tokenProvider := sCtx.MustGet(cfg.JwtConfig.JwtPrefix).(tokenPlugin.TokenProvider)

	str := authStorage.GetNewMySQLConnection(database)
	// init services
	bz := authBusiness.NewAuthenticationBusiness(str, tokenProvider, 60*60*24*30)
	// todo move expire time to config

	return authAPI.NewAuthenticationAPI(bz)
}

func GetAuthenticationService(sCtx serviceCtx.ServiceContext) authBusiness.AuthenticationBusiness {

	database := sCtx.MustGet(common.DbMainName).(*gorm.DB)
	tokenProvider := sCtx.MustGet(cfg.JwtConfig.JwtPrefix).(tokenPlugin.TokenProvider)

	str := authStorage.GetNewMySQLConnection(database)
	// init services
	return authBusiness.NewAuthenticationBusiness(str, tokenProvider, 60*60*24*30)
}

func GetAuthenticationGrpcAPIService(sCtx serviceCtx.ServiceContext) transport.AuthenticationRpcAPI {

	conn, err := grpc.NewClient(
		"localhost:"+grpcConfig.AuthenticationPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to connect to grpc server:", err)
	}

	client := contract.NewAuthenticationServiceClient(conn)
	rpcClient := storage.NewAuthClientGrpc(client)

	bz := authGprcBusiness.NewAuthenticationBusinessGrpc(rpcClient)

	return transport.NewAuthenticationAPI(bz)
}

func GetUserAPIService(sCtx serviceCtx.ServiceContext) userApi.UserAPI {

	database := sCtx.MustGet(common.DbMainName).(*gorm.DB)

	str := userStorage.GetNewMySQLConnection(database)
	bz := userBusiness.NewUserBusiness(str)
	return userApi.NewUserAPI(bz)
}

func GetTodoReactionService(sCtx serviceCtx.ServiceContext) reactionBusiness.ReactionBusiness {

	database := sCtx.MustGet(common.DbMainName).(*gorm.DB)
	ps := sCtx.MustGet(common.PluginPubSub).(pubsub.PubSub)

	str := reactionStorage.GetNewMySQLConnection(database)

	rabbitMQ := registerservice.NewRabbitMQPublisher("amqp://guest:guest@localhost:5672/", "RabbitMQ_Test")
	return reactionBusiness.NewReactionBusiness(str, ps, rabbitMQ)
}

func GetTodoReactionAPIService(sCtx serviceCtx.ServiceContext) reactionAPI.ReactionTodoAPI {

	database := sCtx.MustGet(common.DbMainName).(*gorm.DB)
	ps := sCtx.MustGet(common.PluginPubSub).(pubsub.PubSub)

	str := reactionStorage.GetNewMySQLConnection(database)

	rabbitMQ := registerservice.NewRabbitMQPublisher("amqp://guest:guest@localhost:5672/", "RabbitMQ_Test")
	bz := reactionBusiness.NewReactionBusiness(str, ps, rabbitMQ)

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
