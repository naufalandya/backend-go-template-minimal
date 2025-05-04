package controllers

import (
	"fmt"
	"modular_monolith/server/api/legal/functions"
	"modular_monolith/server/api/legal/utils"
	global "modular_monolith/server/models"

	"github.com/gofiber/fiber/v2"
)

func UploadFile(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(global.Apiresponse{
			Message: "File is required",
			Code:    fiber.StatusBadRequest,
			Status:  false,
			Data:    nil,
		})
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(global.Apiresponse{
			Message: fmt.Sprintf("Failed to open file: %v", err),
			Code:    fiber.StatusInternalServerError,
			Status:  false,
			Data:    nil,
		})
	}
	defer file.Close()

	fileType := utils.GetFileType(fileHeader.Filename)

	var text string
	switch fileType {
	case "pdf":
		text, err = functions.ExtractTextFromPDF(file)
	case "docx":
		text, err = functions.ExtractTextFromDocx(file)
	default:
		return c.Status(fiber.StatusUnsupportedMediaType).JSON(global.Apiresponse{
			Message: "Unsupported file type",
			Code:    fiber.StatusUnsupportedMediaType,
			Status:  false,
			Data:    nil,
		})
	}

	fmt.Println(text)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(global.Apiresponse{
			Message: fmt.Sprintf("Failed to extract text: %v", err),
			Code:    fiber.StatusInternalServerError,
			Status:  false,
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(global.Apiresponse{
		Message: "Text extracted successfully",
		Code:    fiber.StatusOK,
		Status:  true,
		Data:    text,
	})
}
