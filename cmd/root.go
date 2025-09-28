package cmd

import (
	"fmt"
	ssu "github.com/coderconquerer/social-todo/cmd/servicesetup"
	"github.com/coderconquerer/social-todo/pkg/logger"
	"go.uber.org/zap"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Start social todo service",
	Run: func(cmd *cobra.Command, args []string) {
		log := logger.Init()
		defer logger.Sync()
		log = log.Named("service-todo-logger") // gives component name in logs

		service := ssu.NewServices()

		if err := service.Init(); err != nil {
			log.Fatal("error when initializing service: ", zap.Error(err))
		}

		// --- setup components ---
		ssu.StartGrpcServer(service, log)
		ssu.StartHttpServer(service)
		ssu.SetupTracing(log)
		ssu.StartSubscribers(service)
		ssu.StartRabbitMQConsumer(log)

		// --- start service lifecycle ---
		if err := service.Start(); err != nil {
			log.Fatal("error when starting service", zap.Error(err))
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
