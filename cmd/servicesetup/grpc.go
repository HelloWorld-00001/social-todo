package servicesetup

import (
	"go.uber.org/zap"
	"net"

	goService "github.com/200Lab-Education/go-sdk"
	"github.com/coderconquerer/social-todo/composer"
	"github.com/coderconquerer/social-todo/grpc/contract"
	"github.com/coderconquerer/social-todo/module/todotaskreaction/transport/rpc"
	"google.golang.org/grpc"
)

func StartGrpcServer(service goService.Service, log *zap.Logger) {
	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatal("Failed to listen", zap.Error(err))
	}

	grpcServer := grpc.NewServer()
	contract.RegisterItemReactServiceServer(grpcServer,
		rpc.NewRpcService(composer.GetTodoReactionService(service)),
	)

	go func() {
		log.Info("Serving gRPC on 0.0.0.0:8082")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("Failed to listen", zap.Error(err))
		}
	}()
}
