package controllers

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/villagerv/go_inverse/utils"
)

func FetchURL(c *fiber.Ctx) error {
	encodedURL := c.Params("encodedURL")
	decodedURL, err := url.QueryUnescape(encodedURL)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	originalJSON, processedJSON, err := utils.FetchAndProcessURL(decodedURL)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"original":  originalJSON,
		"processed": processedJSON,
	})
}
