package BusinessUseCases

import (
	"bytes"
	"database/sql/driver"
	"errors"
	models2 "github.com/coderconquerer/go-login-app/internal/TodoItem/models"
	"github.com/coderconquerer/go-login-app/internal/common"
	"github.com/coderconquerer/go-login-app/internal/components/uploadProvider"
	"github.com/coderconquerer/go-login-app/internal/user/models"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	img "image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type UploadImageStorage interface {
	UploadImageForTodo(c *gin.Context, id int, image driver.Value) error
}

type UploadAvatarStorage interface {
	UploadUserAvatar(c *gin.Context, id int, image driver.Value) error
}

type UploadImageLogic interface {
	UploadImage(c *gin.Context, image common.Image, destination string, owner common.Entity, id int) (*common.Image, *common.AppError)
}

type UploadFileLogic struct {
	uploadProvider uploadProvider.UploadProvider
	userStore      UploadAvatarStorage
	todoStore      UploadImageStorage
}

func GetNewUploadFileLogic(store UploadImageStorage, userStore UploadAvatarStorage, uploadProvider uploadProvider.UploadProvider) *UploadFileLogic {
	return &UploadFileLogic{
		todoStore:      store,
		userStore:      userStore,
		uploadProvider: uploadProvider,
	}
}

func GetNewUploadFileLogicTemp(uploadProvider uploadProvider.UploadProvider) *UploadFileLogic {
	return &UploadFileLogic{
		uploadProvider: uploadProvider,
	}
}
func (bz *UploadFileLogic) UploadImage(c *gin.Context, image []byte, destination string, owner common.Entity, id int) (*common.Image, *common.AppError) {
	ok, mimeType := isSupportedImage(image)
	if !ok {
		return nil, common.NewBadRequestResponse("unsupported image type: " + mimeType)
	}
	// Decode image config to get width and height
	imgCfg, _, err := img.DecodeConfig(bytes.NewReader(image))
	if err != nil {
		return nil, common.NewBadRequestResponse("unable to decode image metadata")
	}
	// upload to cloud service
	uploadedImg, err2 := bz.uploadProvider.SaveFileUpload(c.Request.Context(), image, destination)
	if err2 != nil {
		return nil, common.NewInternalSeverErrorResponse(err, "error when uploading image", err2.Error())
	}

	// update image metadata after being uploaded
	uploadedImg.Width = imgCfg.Width
	uploadedImg.Height = imgCfg.Height
	uploadedImg.FileSize = len(image)

	err = bz.saveImgToDb(c, uploadedImg, owner, id)
	if err != nil {
		return nil, common.NewInternalSeverErrorResponse(err, "error when uploading image", err.Error())
	}
	return uploadedImg, nil
}

func (bz *UploadFileLogic) saveImgToDb(c *gin.Context, image *common.Image, owner common.Entity, id int) error {
	imgJson, err := image.Value()
	if err != nil {
		return err
	}
	switch owner.ToString() {
	case models.User{}.TableName():
		return bz.userStore.UploadUserAvatar(c, id, imgJson)
	case models2.Todo{}.TableName():
		return bz.todoStore.UploadImageForTodo(c, id, imgJson)
	default:
		return errors.New("invalid owner")
	}
}

func isSupportedImage(data []byte) (bool, string) {
	typeImg := mimetype.Detect(data).String()

	switch typeImg {
	case "image/jpeg", "image/png", "image/gif", "image/webp", "image/bmp", "image/tiff", "image/svg+xml":
		return true, typeImg
	default:
		return false, typeImg
	}
}
