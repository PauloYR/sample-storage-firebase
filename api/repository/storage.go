package repository

import "mime/multipart"

type Storage interface {
	UploadFile(file *multipart.FileHeader) error
	ListFiles(folderName string) (listFiles []string, err error)
}
