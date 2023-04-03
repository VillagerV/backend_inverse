package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func FetchAndProcessURL(url string) ([]map[string]interface{}, []map[string]interface{}, error) {
	originalJSON, err := fetchJSON(url)
	if err != nil {
		return nil, nil, err
	}

	processedData := processJSON(originalJSON)
	processedJSON, ok := processedData.([]map[string]interface{})
	if !ok {
		return nil, nil, errors.New("unable to process JSON data")
	}

	return originalJSON, processedJSON, nil
}

func fetchJSON(url string) ([]map[string]interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unable to fetch data from URL")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jsonData map[string]interface{}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		return nil, err
	}

	results, ok := jsonData["results"].([]interface{})
	if !ok {
		return nil, errors.New("unable to extract 'results' from JSON")
	}

	resultsArray := make([]map[string]interface{}, len(results))
	for i, result := range results {
		resultMap, ok := result.(map[string]interface{})
		if !ok {
			return nil, errors.New("unable to convert 'results' element to map")
		}
		resultsArray[i] = resultMap
	}

	return resultsArray, nil
}

func processJSON(data interface{}) interface{} {
	switch v := data.(type) {
	case []map[string]interface{}:
		processedData := make([]map[string]interface{}, len(v))

		for i, item := range v {
			processedData[i] = processJSON(item).(map[string]interface{})
		}

		return processedData
	case map[string]interface{}:
		processedData := make(map[string]interface{})

		for key, value := range v {
			processedKey := reverseString(key)

			switch value.(type) {
			case string:
				processedData[processedKey] = reverseString(value.(string))
			case []interface{}:
				processedData[processedKey] = reverseList(value.([]interface{}))
			case map[string]interface{}:
				processedData[processedKey] = processJSON(value.(map[string]interface{}))
			default:
				processedData[processedKey] = value
			}
		}

		return processedData
	default:
		return data
	}
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func reverseList(l []interface{}) []interface{} {
	length := len(l)
	reversed := make([]interface{}, length)
	for i, item := range l {
		reversed[length-1-i] = item
	}
	return reversed
}
