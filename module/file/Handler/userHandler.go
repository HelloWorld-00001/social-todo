package Handler

import "github.com/coderconquerer/social-todo/module/file/BusinessUseCases"

type UploadHandler struct {
	UploadBz *BusinessUseCases.UploadFileLogic
}

func NewUploadHandler(GetUserBz *BusinessUseCases.UploadFileLogic) *UploadHandler {
	return &UploadHandler{
		UploadBz: GetUserBz,
	}
}
