package handler

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/suraboy/upload-file-worksho/app/core/model"
	"github.com/suraboy/upload-file-worksho/app/core/service"
	"io/ioutil"
	"net/http"
)

type fileWorkShopHandler struct {
	service service.FileService
}

func SetupProcess(service service.FileService) FileWorkShopHandler {
	return &fileWorkShopHandler{
		service: service,
	}
}

type FileWorkShopHandler interface {
	UploadFilesHandler(c *fiber.Ctx) error
	GetFilesHandler(c *fiber.Ctx) error
}

func (hdl *fileWorkShopHandler) UploadFilesHandler(c *fiber.Ctx) error {
	ctxLog := context.Background()
	// Parse the multipart form data to get the file
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Failed to parse form file")
	}

	// Open the file
	fileContent, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to open file")
	}
	defer fileContent.Close()

	// Read the file data
	fileData := make([]byte, file.Size)
	_, err = fileContent.Read(fileData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to read file data")
	}

	url, err := hdl.service.UploadFileService(ctxLog, file.Filename, fileData)

	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to upload file")
	}

	res := model.FileListResponse{
		Url: "http:localhost:3000/" + url,
	}

	return c.JSON(res)
}

func (hdl *fileWorkShopHandler) GetFilesHandler(c *fiber.Ctx) error {
	objectName := c.Params("file_url")

	url, err := hdl.service.GetFileSignedURLService(objectName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to open file")
	}

	// Make an HTTP GET request to the signed URL to fetch the file contents
	resp, err := http.Get(url)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to open file")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.Status(http.StatusInternalServerError).SendString("Fail to open file")
	}
	// Read the file contents
	fileContents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to read file contents")
	}

	// Set the appropriate content type based on the file type
	contentType := resp.Header.Get("Content-Type")
	c.Set("Content-Type", contentType)

	// Return the file contents in the response
	return c.Send(fileContents)
}
