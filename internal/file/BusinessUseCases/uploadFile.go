package BusinessUseCases

import (
	"github.com/coderconquerer/go-login-app/internal/common"
	"github.com/coderconquerer/go-login-app/internal/components/uploadProvider"
	"github.com/gin-gonic/gin"
)

type UploadImageStorage interface {
	SaveImageMetadata(c *gin.Context, image common.Image) (*common.Image, error)
}

type UploadImageLogic interface {
	UploadImage(c *gin.Context, image common.Image, destination string) (*common.Image, *common.AppError)
}

type UploadFileLogic struct {
	uploadProvider uploadProvider.UploadProvider
	store          UploadImageStorage
}

func GetNewUploadFileLogic(store UploadImageStorage, uploadProvider uploadProvider.UploadProvider) *UploadFileLogic {
	return &UploadFileLogic{
		store:          store,
		uploadProvider: uploadProvider,
	}
}

func GetNewUploadFileLogicTemp(uploadProvider uploadProvider.UploadProvider) *UploadFileLogic {
	return &UploadFileLogic{
		uploadProvider: uploadProvider,
	}
}
func (bz *UploadFileLogic) UploadImage(c *gin.Context, image []byte, destination string) (*common.Image, *common.AppError) {
	res, err := bz.uploadProvider.SaveFileUpload(c.Request.Context(), image, destination)
	if err != nil {
		return nil, common.NewInternalSeverErrorResponse(err, "error when uploading image", err.Error())
	}
	return res, nil
}
