package servicesetup

import (
	goService "github.com/200Lab-Education/go-sdk"
	"github.com/coderconquerer/social-todo/subscribers"
)

func StartSubscribers(service goService.Service) {
	_ = subscribers.NewEngine(service).Start()
}
