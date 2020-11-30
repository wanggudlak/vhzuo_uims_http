package requests

import "mime/multipart"

type FileUploadRequest struct {
	SPMfullCode  string                `form:"spm"`
	UploadedFile *multipart.FileHeader `form:"file"`
}

type UploadRequest struct {
	Key          string                `form:"key"`
	ClientID     int                   `form:"client_id"`
	UploadedFile *multipart.FileHeader `form:"file"`
}
