package servicesetup

import (
	"github.com/coderconquerer/social-todo/configs"
	rpc2 "github.com/coderconquerer/social-todo/module/authentication/transport/rpc"
	"go.uber.org/zap"
	"net"

	goService "github.com/200Lab-Education/go-sdk"
	"github.com/coderconquerer/social-todo/composer"
	"github.com/coderconquerer/social-todo/grpc/contract"
	"github.com/coderconquerer/social-todo/module/todotaskreaction/transport/rpc"
	"google.golang.org/grpc"
)

var grpcConfig = configs.LoadGrpcPort()

func StartTodoReactionGrpcServer(service goService.Service, log *zap.Logger) {
	if grpcConfig.GrpcStart {
		log.Info("skip start TodoReactionGrpc grpc")
		return
	}
	lis, err := net.Listen("tcp", ":"+grpcConfig.TodoReactionPort)
	if err != nil {
		log.Fatal("Failed to listen", zap.Error(err))
	}

	grpcServer := grpc.NewServer()
	contract.RegisterItemReactServiceServer(grpcServer,
		rpc.NewRpcService(composer.GetTodoReactionService(service)),
	)

	go func() {
		log.Info("Serving TodoReactionGrpc", zap.String("addr", "0.0.0.0:"+grpcConfig.AuthenticationPort))
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("Failed to listen", zap.Error(err))
		}
	}()
}

func StartAuthenticationGrpcServer(service goService.Service, log *zap.Logger) {
	if grpcConfig.GrpcStart {
		log.Info("skip start AuthenticationGrpc grpc")
		return
	}
	lis, err := net.Listen("tcp", ":"+grpcConfig.AuthenticationPort)
	if err != nil {
		log.Fatal("Failed to listen", zap.Error(err))
	}

	grpcServer := grpc.NewServer()
	contract.RegisterAuthenticationServiceServer(grpcServer,
		rpc2.NewRPCServer(composer.GetAuthenticationService(service)),
	)

	go func() {
		log.Info("Serving AuthenticationGrpc", zap.String("addr", "0.0.0.0:"+grpcConfig.AuthenticationPort))
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("Failed to listen", zap.Error(err))
		}
	}()
}
