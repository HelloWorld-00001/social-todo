package plugin

import (
	"context"
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/plugin"
)

type UploadProvider interface {
	SaveFileUpload(ctx context.Context, data []byte, destination string) (*common.Image, error)
	plugin.PluginBase
}
