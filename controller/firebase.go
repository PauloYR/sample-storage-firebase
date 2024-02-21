package controller

import (
	"mime/multipart"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

func (controller Controller) Upload(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["files"]

	var wg sync.WaitGroup
	sizeFiles := len(files)
	errChan := make(chan error, sizeFiles)
	for _, file := range files {
		wg.Add(1)
		go func(wg *sync.WaitGroup, errChan chan<- error, file *multipart.FileHeader) {
			defer wg.Done()
			errChan <- controller.storage.UploadFile(file)
		}(&wg, errChan, file)
	}

	wg.Wait()
	close(errChan)
	for result := range errChan {
		if result != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Upload successful!"})
}

func (controller Controller) ListFilesInFolder(ctx *gin.Context) {
	files, err := controller.storage.ListFiles("")

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, files)
}
