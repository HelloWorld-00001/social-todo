package business

import (
	"bytes"
	"context"
	"database/sql/driver"
	"errors"
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/module/file/entity"
	todoEntity "github.com/coderconquerer/social-todo/module/todo/entity"
	userEntity "github.com/coderconquerer/social-todo/module/user/entity"
	uploadProvider "github.com/coderconquerer/social-todo/plugin/uploadProvider"
	"github.com/gabriel-vasile/mimetype"
	img "image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type UploadTodoImageStorage interface {
	UploadImageForTodo(c context.Context, id int, image driver.Value) error
}

type UploadUserImageStorage interface {
	UploadUserAvatar(c context.Context, id int, image driver.Value) error
}

type UploadImageBusiness interface {
	UploadImage(c context.Context, image []byte, destination string, owner common.Entity, id int) (*common.Image, error)
}

type uploadFileBusiness struct {
	uploadProvider         uploadProvider.UploadProvider
	uploadTodoImageStorage UploadTodoImageStorage
	uploadUserImageStorage UploadUserImageStorage
}

func NewUploadImageLogic(uploadTodoImageStorage UploadTodoImageStorage,
	uploadUserImageStorage UploadUserImageStorage,
	uploadProvider uploadProvider.UploadProvider) UploadImageBusiness {
	return &uploadFileBusiness{
		uploadTodoImageStorage: uploadTodoImageStorage,
		uploadUserImageStorage: uploadUserImageStorage,
		uploadProvider:         uploadProvider,
	}
}

func (bz *uploadFileBusiness) UploadImage(c context.Context, image []byte, destination string, owner common.Entity, id int) (*common.Image, error) {
	ok, mimeType := isSupportedImage(image)
	if !ok {
		return nil, common.BadRequest.WithError(errors.New("unsupported image type: " + mimeType))
	}
	// Decode image config to get width and height
	imgCfg, _, err := img.DecodeConfig(bytes.NewReader(image))
	if err != nil {
		return nil, common.BadRequest.WithError(errors.New("unable to decode image metadata"))
	}
	// upload to cloud service
	uploadedImg, errUpload := bz.uploadProvider.SaveFileUpload(c, image, destination)
	if errUpload != nil {
		return nil, common.InternalServerError.WithError(entity.ErrCannotUploadFile).WithRootCause(errUpload)
	}

	// update image metadata after being uploaded
	uploadedImg.Width = imgCfg.Width
	uploadedImg.Height = imgCfg.Height
	uploadedImg.FileSize = len(image)

	errBusiness := bz.saveImageMetadataToDb(c, uploadedImg, owner, id)
	if errBusiness != nil {
		return nil, common.InternalServerError.WithError(entity.ErrCannotUploadFile).WithRootCause(errBusiness)
	}
	return uploadedImg, nil
}

func (bz *uploadFileBusiness) saveImageMetadataToDb(c context.Context, image *common.Image, owner common.Entity, id int) error {
	imgJson, err := image.Value()
	if err != nil {
		return err
	}
	switch owner.ToString() {
	case userEntity.User{}.TableName():
		return bz.uploadUserImageStorage.UploadUserAvatar(c, id, imgJson)
	case todoEntity.Todo{}.TableName():
		return bz.uploadTodoImageStorage.UploadImageForTodo(c, id, imgJson)
	default:
		return common.BadRequest.WithError(errors.New("unknown owner type" + owner.ToString()))
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
