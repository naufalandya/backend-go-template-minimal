package controllers

import (
	"fmt"
	"io"
	"log"
	global "modular_monolith/models"
	pbapi "modular_monolith/protobuf/api"
	"modular_monolith/server/api/user/models"
	"modular_monolith/server/api/user/provider"
	"modular_monolith/server/api/user/services"
	"modular_monolith/server/client"
	"modular_monolith/server/functions"
	"strings"

	fiber "github.com/gofiber/fiber/v2"
)

func GetAllUsers(c *fiber.Ctx) error {
	users, err := provider.FetchAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(global.Apiresponse{
			Code:    fiber.StatusInternalServerError,
			Status:  false,
			Message: fmt.Sprintf("Failed to fetch users : %s", err.Error()),
		})
	}
	return c.Status(fiber.StatusOK).JSON(global.Apiresponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "Success",
		Data:    users,
	})
}

func GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := provider.FetchUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(global.Apiresponse{
			Code:    fiber.StatusOK,
			Status:  false,
			Message: fmt.Sprintf("Failed to fetch users : %s", err.Error()),
		})
	}
	return c.Status(fiber.StatusOK).JSON(global.Apiresponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "Success",
		Data:    user,
	})
}

func RegisterUserSimple(c *fiber.Ctx) error {
	var input models.RegisterRequest

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(global.Apiresponse{
			Code:    fiber.StatusBadRequest,
			Status:  false,
			Message: fmt.Sprintf("Cannot parse your body~ (´；ω；｀) : %s", err.Error()),
		})
	}

	if err := functions.FuckOffHackerByJSON(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(global.Apiresponse{
			Code:    fiber.StatusBadRequest,
			Status:  false,
			Message: fmt.Sprintf("Suspicious input detected~ (｀_´) : %s", err.Error()),
		})
	}

	if errs := functions.ValidateStruct(input); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(global.Apiresponse{
			Code:    fiber.StatusBadRequest,
			Status:  false,
			Message: fmt.Sprintf("Validation failed~ (｀_´) : %s", errs[0]),
		})
	}

	// exist, err := provider.IsEmailOrUsernameExist(input.Email, input.Username)
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(global.Apiresponse{
	// 		Code:    fiber.StatusInternalServerError,
	// 		Status:  false,
	// 		Message: "Oops, something went wrong and its on us~ (｡•́︿•̀｡)",
	// 	})
	// }
	// if exist {
	// 	return c.Status(fiber.StatusBadRequest).JSON(global.Apiresponse{
	// 		Code:    fiber.StatusBadRequest,
	// 		Status:  false,
	// 		Message: "Email or Username already used~ (｀_´)",
	// 	})
	// }

	hashedPassword, err := functions.HashPassword(input.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(global.Apiresponse{
			Code:    fiber.StatusInternalServerError,
			Status:  false,
			Message: "Oops, something went wrong and its on us~ (｡•́︿•̀｡)",
		})
	}

	fmt.Println(hashedPassword)

	// if err := provider.CreateUser(hashedPassword, input); err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(global.Apiresponse{
	// 		Code:    fiber.StatusInternalServerError,
	// 		Status:  false,
	// 		Message: "Cannot create user~ (つω≦ )",
	// 	})
	// }

	if err := services.PublishMessage(input.FullName); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(global.Apiresponse{
			Code:    fiber.StatusInternalServerError,
			Status:  false,
			Message: "Cannot publish message to Redis ~ (つω≦ )",
		})
	}

	if err := services.PublishMessageRabbit("your-queue-name", input.FullName); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(global.Apiresponse{
			Code:    fiber.StatusInternalServerError,
			Status:  false,
			Message: "Cannot publish message to Rabbitt ~ (つω≦ )",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(global.Apiresponse{
		Code:    fiber.StatusCreated,
		Status:  true,
		Message: "User created successfully~ (๑˃̵ᴗ˂̵)و",
	})
}

func SayHello(c *fiber.Ctx) error {
	var input models.RegisterRequest

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(global.Apiresponse{
			Code:    fiber.StatusBadRequest,
			Status:  false,
			Message: fmt.Sprintf("Cannot parse your body~ (´；ω；｀) : %s", err.Error()),
		})
	}

	res, err := client.Clients.HelloWorldClient.SayHello(c.Context(), &pbapi.HelloRequest{
		Name: input.FullName, // or hardcoded "Andya"
	})

	if err != nil {
		log.Fatalf("Error calling gRPC service: %v", err)
	}
	log.Printf("Response: %v", res)

	return c.Status(fiber.StatusOK).JSON(global.Apiresponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "Okiie dokkie ~ (๑˃̵ᴗ˂̵)و",
	})
}

// Constants for allowed file types and maximum file size
const maxFileSize = 10 * 1024 * 1024 // 10MB

func UploadFile(c *fiber.Ctx) error {
	// Parse the form data to get the file
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(global.Apiresponse{
			Code:    fiber.StatusBadRequest,
			Status:  false,
			Message: fmt.Sprintf("Failed to upload file: %s", err.Error()),
		})
	}

	// Check the file size
	if file.Size > maxFileSize {
		return c.Status(fiber.StatusBadRequest).JSON(global.Apiresponse{
			Code:    fiber.StatusBadRequest,
			Status:  false,
			Message: "File is too large. Maximum allowed size is 10MB.",
		})
	}

	validExtensions := []string{".jpg", ".jpeg", ".png", ".pdf"}
	fileExt := strings.ToLower(file.Filename[strings.LastIndex(file.Filename, "."):])

	if !functions.Contains(validExtensions, fileExt) {
		return c.Status(fiber.StatusBadRequest).JSON(global.Apiresponse{
			Code:    fiber.StatusBadRequest,
			Status:  false,
			Message: "Invalid file type. Only .jpg, .jpeg, .png, and .pdf are allowed.",
		})
	}

	fileHeader, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(global.Apiresponse{
			Code:    fiber.StatusInternalServerError,
			Status:  false,
			Message: "Error reading file header.",
		})
	}
	defer fileHeader.Close()

	if !functions.IsValidFileType(fileHeader) {
		return c.Status(fiber.StatusBadRequest).JSON(global.Apiresponse{
			Code:    fiber.StatusBadRequest,
			Status:  false,
			Message: "File type does not match the expected formats.",
		})
	}

	if err := functions.ScanForVirus(fileHeader); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(global.Apiresponse{
			Code:    fiber.StatusBadRequest,
			Status:  false,
			Message: fmt.Sprintf("File contains a virus or malware: %s", err.Error()),
		})
	}

	fileData, err := io.ReadAll(fileHeader)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(global.Apiresponse{
			Code:    fiber.StatusInternalServerError,
			Status:  false,
			Message: "Error reading file content.",
		})
	}

	return c.Status(fiber.StatusOK).Send(fileData)
}

// func CreateUser(c *fiber.Ctx) error {
// 	var input models.UserInput

// 	if err := c.BodyParser(&input); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": "Cannot parse your body~ (´；ω；｀)",
// 			"error":   err.Error(),
// 		})
// 	}

// 	if sanitizedValue, err := functions.SuperSecureSanitize(input.Name); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": "Suspicious input in 'Name' field~ (｀_´)",
// 			"error":   err.Error(),
// 		})
// 	} else {
// 		input.Name = sanitizedValue.(string)
// 	}

// 	if sanitizedValue, err := functions.SuperSecureSanitize(input.Email); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": "Suspicious input in 'Email' field~ (｀_´)",
// 			"error":   err.Error(),
// 		})
// 	} else {
// 		input.Email = sanitizedValue.(string)
// 	}

// 	if sanitizedValue, err := functions.SuperSecureSanitize(input.Age); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": "Suspicious input in 'Age' field~ (｀_´)",
// 			"error":   err.Error(),
// 		})
// 	} else {
// 		input.Age = sanitizedValue.(int)
// 	}

// 	if errs := functions.ValidateStruct(input); errs != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": "Validation failed~",
// 			"errors":  errs,
// 		})
// 	}

// 	return c.JSON(fiber.Map{
// 		"message": "User created successfully! (っ＾▿＾)💨",
// 		"data":    input,
// 	})
// }
