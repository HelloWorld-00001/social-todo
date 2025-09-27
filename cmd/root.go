package cmd

import (
	"context"
	"contrib.go.opencensus.io/exporter/jaeger"
	"fmt"
	goService "github.com/200Lab-Education/go-sdk"
	"github.com/200Lab-Education/go-sdk/plugin/storage/sdkgorm"
	"github.com/coderconquerer/social-todo/cache"
	"github.com/coderconquerer/social-todo/cmd/registerservice"
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/configs"
	"github.com/coderconquerer/social-todo/docs"
	"github.com/coderconquerer/social-todo/grpc/contract"
	"github.com/coderconquerer/social-todo/middleware"
	accBuc "github.com/coderconquerer/social-todo/module/account/BusinessUseCases"
	accountHdl "github.com/coderconquerer/social-todo/module/account/Handler"
	accStorage "github.com/coderconquerer/social-todo/module/account/Storage"
	"github.com/coderconquerer/social-todo/module/file/BusinessUseCases"
	"github.com/coderconquerer/social-todo/module/file/Handler"
	todoBuc "github.com/coderconquerer/social-todo/module/todoItem/BusinessUseCases"
	todoHdl "github.com/coderconquerer/social-todo/module/todoItem/Handler"
	"github.com/coderconquerer/social-todo/module/todoItem/Handler/rpc"
	"github.com/coderconquerer/social-todo/module/todoItem/Repository"
	todoStorage "github.com/coderconquerer/social-todo/module/todoItem/Storage"
	restapi2 "github.com/coderconquerer/social-todo/module/todoItem/Storage/grpc"
	"github.com/coderconquerer/social-todo/module/todoItem/Storage/restapi"
	userBuc "github.com/coderconquerer/social-todo/module/user/BusinessUseCases"
	userHdl "github.com/coderconquerer/social-todo/module/user/Handler"
	userStorage "github.com/coderconquerer/social-todo/module/user/Storage"
	BusinessUseCases2 "github.com/coderconquerer/social-todo/module/userReactItem/BusinessUseCases"
	Handler2 "github.com/coderconquerer/social-todo/module/userReactItem/Handler"
	rpc2 "github.com/coderconquerer/social-todo/module/userReactItem/Handler/rpc"
	"github.com/coderconquerer/social-todo/module/userReactItem/Storage"
	"github.com/coderconquerer/social-todo/plugin/redis"
	"github.com/coderconquerer/social-todo/plugin/rpc_caller"
	tokenPlugin "github.com/coderconquerer/social-todo/plugin/tokenProviders"
	"github.com/coderconquerer/social-todo/plugin/tokenProviders/jwtProvider"
	uploadPlugin "github.com/coderconquerer/social-todo/plugin/uploadProvider"
	"github.com/coderconquerer/social-todo/plugin/uploadProvider/s3provider"
	"github.com/coderconquerer/social-todo/pubsub"
	"github.com/coderconquerer/social-todo/subscribers"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/cobra"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
	"strconv"
	"time"
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
		goService.WithInitRunnable(pubsub.NewLocalPubsub(common.PluginPubSub)),
		goService.WithInitRunnable(rpc_caller.NewRpcCaller(common.PluginRPC)),
		goService.WithInitRunnable(redis.NewRedisDB("redisDb", common.PluginRedis)),
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
		// Create a listener on TCP port
		lis, err := net.Listen("tcp", ":8082")
		if err != nil {
			log.Fatalln("Failed to listen:", err)
		}

		// Create a gRPC server object
		s := grpc.NewServer()
		service.HTTPServer().AddHandler(func(engine *gin.Engine) {
			database := service.MustGet(common.DbMainName).(*gorm.DB)
			tokenProvider := service.MustGet(cfg.JwtConfig.JwtPrefix).(tokenPlugin.TokenProvider)
			s3Provider := service.MustGet(awsCfg.S3Prefix).(uploadPlugin.UploadProvider)
			ps := service.MustGet(common.PluginPubSub).(pubsub.PubSub)
			rpcCaller := service.MustGet(common.PluginRPC).(rpc_caller.RpcCaller)

			// Create a client connection to the gRPC server we just started
			// This is where the gRPC-Gateway proxies the requests

			conn, err := grpc.NewClient(
				"0.0.0.0:8082",
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)

			//	TODO implement gateway
			// Create client stub
			client := contract.NewItemReactServiceClient(conn)
			if err != nil {
				log.Fatalln("Failed to dial server:", err)
			}

			// init services
			accStore := accStorage.GetNewMySQLConnection(database)
			registerBsn := accBuc.GetNewRegisterAccountLogic(accStore)
			loginBsn := accBuc.GetNewLoginLogic(accStore, tokenProvider, 60*60*24*30)
			disableBsn := accBuc.GetNewDisableAccountLogic(accStore)
			accountHandler := accountHdl.NewAccountHandler(loginBsn, registerBsn, disableBsn)

			// reaction service
			reactionStore := Storage.GetNewMySQLConnection(database)

			test := registerservice.NewRabbitMQPublisher("amqp://guest:guest@localhost:5672/", "RabbitMQ_Test")
			reactBz := BusinessUseCases2.GetNewReactTodoItemLogic(reactionStore, ps, test)
			unReactBz := BusinessUseCases2.GetNewUnreactTodoItemLogic(reactionStore, ps)
			listRUBz := BusinessUseCases2.GetNewGetListReactedUsersLogic(reactionStore)
			reactHandler := Handler2.NewReactionHandler(reactBz, unReactBz, listRUBz)
			_ = restapi.NewTodoReactService(rpcCaller.GetServiceUrl(), rpcCaller.GetClient())
			rpcClient := restapi2.NewRpcClient(client)
			// to do service
			todoStore := todoStorage.GetNewMySQLConnection(database)
			getTodoDetailBz := todoBuc.GetNewGetTodoDetailLogic(todoStore)
			getTodoListRepo := Repository.GetNewTodoListWithReactRepo(todoStore, rpcClient)
			getTodoListBz := todoBuc.GetNewGetTodoListLogic(getTodoListRepo)
			createTodoListBz := todoBuc.GetNewCreateTodoLogic(todoStore)
			deleteTodoListBz := todoBuc.GetNewDeleteTodoItemLogic(todoStore)
			todoHandler := todoHdl.NewTodoHandler(getTodoDetailBz, getTodoListBz, createTodoListBz, deleteTodoListBz)

			// rpc service
			getTotalReactBz := todoBuc.GetNewGetTodoItemTotalReactLogic(reactionStore)
			todoRpcHandler := rpc.NewRpcHandler(getTotalReactBz)
			// user service
			userStore := userStorage.GetNewMySQLConnection(database)
			getUserProfileBz := userBuc.GetNewFindUserLogic(userStore)
			userHandler := userHdl.NewUserHandler(getUserProfileBz)

			// aws services
			uploadBsn := BusinessUseCases.GetNewUploadFileLogic(todoStore, userStore, s3Provider)
			uploadHandler := Handler.NewUploadHandler(uploadBsn)
			swaggerSetup()

			// redis cache
			redisService := cache.NewRedisCache(service)
			userCache := cache.NewUserCaching(redisService, accStore)
			// init middleware
			authAdmin := middleware.RequireAuth(tokenProvider, userCache, common.AdminRole.ToString())
			authUser := middleware.RequireAuth(tokenProvider, userCache, []string{common.AdminRole.ToString(), common.UserRole.ToString()}...) // multi-role access

			// rpc
			// Attach the Greeter service to the server
			contract.RegisterItemReactServiceServer(s, rpc2.NewRpcService(listRUBz))

			// Serve gRPC Server
			go func() {
				log.Println("Serving gRPC on 0.0.0.0:8082")
				log.Fatal(s.Serve(lis))
			}()
			//engine.Use(gin.Logger())
			//engine.Use(middleware.CustomRecovery()) // custom middleware
			//auth := engine.Group("/")b
			store := cookie.NewStore([]byte("super-secret-key"))

			v1 := engine.Group("/v1/api")
			{
				todoRoutes := v1.Group("/todo")
				{
					todoRoutes.GET("", authUser, todoHandler.GetToDoList())
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

				auth := v1.Group("/auth")
				{
					auth.POST("/login", accountHandler.Login())
					auth.POST("/register", accountHandler.RegisterAccount())
				}

				v1.GET("/profile", authUser, userHandler.GetUserProfile())
				v1.POST("/disable", authAdmin, accountHandler.DisableAccount())
				v1.POST("/upload", authUser, uploadHandler.UploadImage())

				rpc := v1.Group("/rpc")
				{
					rpc.POST("/get_todo_total_react", todoRpcHandler.GetTodoTotalReact())
				}

			}
			engine.Use(sessions.Sessions("mysession", store))
			engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
			engine.GET("/ping", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "pong",
				})
			})
		})

		// todo-tracer: migrate to OpenTelementry
		jg, err := jaeger.NewExporter(jaeger.Options{
			AgentEndpoint: "localhost:6831",
			Process: jaeger.Process{
				ServiceName: "social-todo-app",
			}})

		go processFib()
		if err != nil {
			log.Fatalln(err)
		}
		trace.RegisterExporter(jg)
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.ProbabilitySampler(1)})

		_ = subscribers.NewEngine(service).Start()
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

func fib(n int) int {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return fib(n-1) + fib(n-2)
	}
}

func processFib() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	log.Printf("Failed to connect to RabbitMQ rp %v", err)
	fmt.Printf("Failed to connect to RabbitMQ server: %v", err)
	defer conn.Close()

	ch, err := conn.Channel()
	fmt.Println(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"rpc_queue", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	fmt.Println(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	fmt.Println(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	fmt.Println(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		for d := range msgs {
			n, err := strconv.Atoi(string(d.Body))
			fmt.Println(err, "Failed to convert body to integer")

			log.Printf(" [.] fib(%d)", n)
			response := fib(n)

			err = ch.PublishWithContext(ctx,
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          []byte(strconv.Itoa(response)),
				})
			fmt.Println(err, "Failed to publish a message")

			d.Ack(false)
		}
	}()

	log.Printf(" [*] Awaiting RPC requests")
	<-forever
}
