package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type ApiResponse struct {
	OriginalResponse  json.RawMessage   `json:"original_response"`
	ProcessedResponse map[string]string `json:"processed_response"`
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func processResponse(originalResponse map[string]interface{}) map[string]string {
	processedResponse := make(map[string]string)

	for key, value := range originalResponse {
		processedKey := reverseString(key)
		processedValue := ""

		if strValue, ok := value.(string); ok {
			processedValue = reverseString(strValue)
		}

		processedResponse[processedKey] = processedValue
	}

	return processedResponse
}

func fetchAndProcessData(c *fiber.Ctx) error {
	searchVal := c.Query("searchVal")
	returnGeom := c.Query("returnGeom")
	getAddrDetails := c.Query("getAddrDetails")
	pageNum := c.Query("pageNum")

	url := fmt.Sprintf("https://developers.onemap.sg/commonapi/search?searchVal=%s&returnGeom=%s&getAddrDetails=%s&pageNum=%s", searchVal, returnGeom, getAddrDetails, pageNum)

	resp, err := http.Get(url)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch data from the provided URL",
		})
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read response body",
		})
	}

	var originalResponse map[string]interface{}
	json.Unmarshal(body, &originalResponse)

	processedResponse := processResponse(originalResponse)

	apiResponse := ApiResponse{
		OriginalResponse:  json.RawMessage(body),
		ProcessedResponse: processedResponse,
	}

	return c.JSON(apiResponse)
}

func main() {
	app := fiber.New()

	app.Get("/api/fetch", fetchAndProcessData)

	app.Listen(":3001")
}
