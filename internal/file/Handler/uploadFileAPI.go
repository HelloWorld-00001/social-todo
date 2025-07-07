package Handler

import (
	"errors"
	"github.com/coderconquerer/go-login-app/internal/common"
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"net/http"
)

func (uh *UploadHandler) UploadImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse multipart form (optional: set memory limit)
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewBadRequestErrorResponse(errors.New("missing upload file"), "file is quire", ""))
			return
		}

		openedFile, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.NewInternalSeverErrorResponse(errors.New("error when open file"), "Cannot open uploaded file", ""))
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
		_, err = openedFile.Read(dataBytes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.NewInternalSeverErrorResponse(errors.New("error when read file"), "Cannot read uploaded file", ""))
			return
		}

		// Optional: clean or generate file name
		destination := "social-todo-list/" + file.Filename
		// Optionally add UUID, timestamp, etc.

		// Upload to S3
		img, err2 := uh.UploadBz.UploadImage(c, dataBytes, destination)
		if err2 != nil {
			c.JSON(http.StatusInternalServerError, common.NewInternalSeverErrorResponse(err2, "cannot upload image to S3 cloud", err2.Error()))
			return
		}

		// Respond with the image info
		c.JSON(http.StatusOK, common.SimpleResponse(img))
	}
}
