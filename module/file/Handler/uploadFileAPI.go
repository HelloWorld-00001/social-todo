package Handler

import (
	"errors"
	"fmt"
	common2 "github.com/coderconquerer/social-todo/common"
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func (uh *UploadHandler) UploadImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse multipart form (optional: set memory limit)
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, common2.NewBadRequestResponseWithError(errors.New("missing upload file"), "file is quire", ""))
			return
		}

		openedFile, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, common2.NewInternalSeverErrorResponse(errors.New("error when open file"), "Cannot open uploaded file", ""))
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
			c.JSON(http.StatusInternalServerError, common2.NewInternalSeverErrorResponse(errors.New("error when read file"), "Cannot read uploaded file", ""))
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
			c.JSON(http.StatusBadRequest, common2.NewBadRequestResponse("owner and owner_id are required"))
			return
		}
		ownerId, err := strconv.Atoi(ownerIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, common2.NewBadRequestResponseWithError(err, "invalid owner_id", "owner_id must be integer"))
			return
		}

		// Upload to S3
		img, err2 := uh.UploadBz.UploadImage(c, dataBytes, destination, common2.EntityFromString(ownerStr), ownerId)
		if err2 != nil {
			c.JSON(http.StatusInternalServerError, common2.NewInternalSeverErrorResponse(err2, "cannot upload image to S3 cloud", err2.Error()))
			return
		}

		// Respond with the image info
		c.JSON(http.StatusOK, common2.SimpleResponse(img))
	}
}
