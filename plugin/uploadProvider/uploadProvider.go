package plugin

import (
	"context"
	"github.com/coderconquerer/go-login-app/common"
	"github.com/coderconquerer/go-login-app/plugin"
)

type UploadProvider interface {
	SaveFileUpload(ctx context.Context, data []byte, destination string) (*common.Image, error)
	plugin.PluginBase
}
