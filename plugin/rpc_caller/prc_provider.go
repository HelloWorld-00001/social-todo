package rpc_caller

import (
	"flag"
	"github.com/coderconquerer/social-todo/plugin"
	"github.com/go-resty/resty/v2"
	"log"
	"os"
	"path/filepath"
)

type RpcCaller interface {
	plugin.PluginBase
	GetServiceUrl() string
	GetClient() *resty.Client
}

type rpcCaller struct {
	name       string
	client     *resty.Client
	serviceUrl string
	logger     *log.Logger
	logPath    string
}

func NewRpcCaller(name string) *rpcCaller {
	return &rpcCaller{name: name}
}

func (r *rpcCaller) Name() string {
	return r.name
}

func (r *rpcCaller) GetPrefix() string {
	return r.name
}

func (r *rpcCaller) Get() interface{} {
	return r
}

func (r *rpcCaller) InitFlags() {
	flag.StringVar(&r.serviceUrl, "rpc-todo-react-service-url", "http://localhost:8081/v1/api/rpc/", "use to connect to todo react service")
	flag.StringVar(&r.logPath, "rpc-log-file-path", "", "use to specific log file path")

	// todo: 200lab.io go sdk is wrong, it doesn't auto call configure,
	// so this is temp fix for it
	_ = r.Configure()
}

func (r *rpcCaller) Configure() error {
	r.client = resty.New()
	logFile, err := getLogFile(r.logPath)
	if err != nil {
		log.Fatal(err)
	}
	r.logger = log.New(logFile, "[rpc_caller] ", log.LstdFlags|log.Lshortfile)
	return nil
}

func (r *rpcCaller) Run() error {
	return nil
}

func (r *rpcCaller) Stop() <-chan bool {
	ch := make(chan bool)
	go func() {
		ch <- true
	}()
	return ch
}

func (r *rpcCaller) GetServiceUrl() string {
	return r.serviceUrl
}

func (r *rpcCaller) GetClient() *resty.Client {
	return r.client
}

func getLogFile(logFolder string) (*os.File, error) {
	logDir := filepath.Join(".", "log") // relative path: ./log

	if logFolder != "" {
		logDir = logFolder
	}

	// Step 2: Ensure the directory exists
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("failed to create log directory: %v, folder %s does not exist ", err, logDir)
		return nil, err
	}

	// Step 3: Build full log file path
	logFilePath := filepath.Join(logDir, "rpc.log")

	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
		return nil, err
	}

	return file, nil
}
