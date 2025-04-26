package controllers

import (
	"fmt"
	global "modular_monolith/model"
	"modular_monolith/module/user/models"
	"modular_monolith/module/user/provider"
	"modular_monolith/server/functions"

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

func CreateUser(c *fiber.Ctx) error {
	var input models.UserInput

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

	// Do Something With Data, Always Clean Code !, Always Put Main Logic Inside Controller, Dont Detail It !

	return c.Status(fiber.StatusCreated).JSON(global.Apiresponse{
		Code:    fiber.StatusCreated,
		Status:  true,
		Message: "Success",
		Data:    input,
	})
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
