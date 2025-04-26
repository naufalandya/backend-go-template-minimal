package controllers

import (
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
			Message: "Failed to fetch users",
			Error:   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "Success",
		"data":    users,
	})
}

func GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := provider.FetchUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "Success",
		"data":    user,
	})
}

func CreateUser(c *fiber.Ctx) error {
	var input models.UserInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot parse your body~ (´；ω；｀)",
			"error":   err.Error(),
		})
	}

	if sanitizedValue, err := functions.SuperSecureSanitize(input.Name); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Suspicious input in 'Name' field~ (｀_´)",
			"error":   err.Error(),
		})
	} else {
		input.Name = sanitizedValue.(string)
	}

	if sanitizedValue, err := functions.SuperSecureSanitize(input.Email); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Suspicious input in 'Email' field~ (｀_´)",
			"error":   err.Error(),
		})
	} else {
		input.Email = sanitizedValue.(string)
	}

	if sanitizedValue, err := functions.SuperSecureSanitize(input.Age); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Suspicious input in 'Age' field~ (｀_´)",
			"error":   err.Error(),
		})
	} else {
		input.Age = sanitizedValue.(int)
	}

	if errs := functions.ValidateStruct(input); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed~",
			"errors":  errs,
		})
	}

	return c.JSON(fiber.Map{
		"message": "User created successfully! (っ＾▿＾)💨",
		"data":    input,
	})
}

func CreateUserV2(c *fiber.Ctx) error {
	var input models.UserInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot parse your body~ (´；ω；｀)",
			"error":   err.Error(),
		})
	}

	if err := functions.AutoSuperSanitizeStruct(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Suspicious input detected~ (｀_´)",
			"error":   err.Error(),
		})
	}

	if errs := functions.ValidateStruct(input); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed~",
			"errors":  errs,
		})
	}

	return c.JSON(fiber.Map{
		"message": "User created successfully! (っ＾▿＾)💨",
		"data":    input,
	})
}
