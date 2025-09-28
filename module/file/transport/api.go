package uploadImageAPI

import (
	"errors"
	"fmt"
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/module/file/business"
	"github.com/coderconquerer/social-todo/module/file/entity"
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type UploadImageAPI interface {
	UploadImage() gin.HandlerFunc
}

type uploadAPI struct {
	UploadBz business.UploadImageBusiness
}

func NewUploadImageAPI(GetUserBz business.UploadImageBusiness) UploadImageAPI {
	return &uploadAPI{
		UploadBz: GetUserBz,
	}
}

func (uh *uploadAPI) UploadImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse multipart form (optional: set memory limit)
		file, err := c.FormFile("file")
		if err != nil {
			common.RespondError(c, common.BadRequest.WithError(entity.ErrMissingUploadFile).WithRootCause(err))
			return
		}

		openedFile, errOpen := file.Open()
		if errOpen != nil {
			common.RespondError(c, common.BadRequest.WithError(entity.ErrCannotOpenFile).WithRootCause(errOpen))
			return
		}
		defer func(openedFile multipart.File) {
			err := openedFile.Close()
			if err != nil {
				log.Fatalln("Cannot close uploaded file", err)
			}
		}(openedFile)

		// Read file bytes
		dataBytes := make([]byte, file.Size)
		_, errRead := openedFile.Read(dataBytes)
		if errRead != nil {
			common.RespondError(c, common.BadRequest.WithError(entity.ErrCannotReadFile).WithRootCause(errRead))
			return
		}

		// Optional: clean or generate file name
		ext := strings.ToLower(filepath.Ext(file.Filename)) // Safe: handles case and ensures lowercase
		newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		destination := "social-todo-list/" + newFileName

		// Optionally add UUID, timestamp, etc.
		// Get other form fields
		ownerStr := c.PostForm("owner")
		ownerIdStr := c.PostForm("owner_id")

		if ownerStr == "" || ownerIdStr == "" {
			common.RespondError(c, common.BadRequest.WithError(errors.New("owner and owner_id are required")))
			return
		}
		ownerId, errConvert := strconv.Atoi(ownerIdStr)
		if errConvert != nil {
			common.RespondError(c, common.BadRequest.WithError(errors.New("owner must be an integer")).WithRootCause(errConvert))
			return
		}

		// Upload to S3
		img, errUpload := uh.UploadBz.UploadImage(c, dataBytes, destination, common.EntityFromString(ownerStr), ownerId)
		if errUpload != nil {
			common.RespondError(c, errUpload)
			return
		}

		// Respond with the image info
		c.JSON(http.StatusOK, common.SimpleResponse(img))
	}
}
