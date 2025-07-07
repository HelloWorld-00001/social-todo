package uploadProvider

import (
	"context"
	"github.com/coderconquerer/go-login-app/internal/common"
)

type UploadProvider interface {
	SaveFileUpload(ctx context.Context, data []byte, destination string) (*common.Image, error)
}
