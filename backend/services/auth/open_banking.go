package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"main/config"
	"net/http"
)

func CreateLeanCustomer(openBankingId string) (string, error) {
	url := "https://sandbox.leantech.me/customers/v1/"
	leanAppToken := config.GetConfig().GetString("lean-app-token")

	// Create the request body
	data := map[string]string{
		"app_user_id": openBankingId,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("lean-app-token", leanAppToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		log.Error().Msgf("Request to Lean failed with status code: %d", resp.StatusCode)
		return "", fmt.Errorf("Request to Lean failed with status code: %d", resp.StatusCode)
	}

	// Read the response body
	var responseBody map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return "", err
	}

	// Extract the app_user_id from the response
	appUserID, ok := responseBody["customer_id"].(string)
	if !ok {
		log.Error().Msgf("customer_id not found in the response")
		return "", fmt.Errorf("customer_id not found in the response")
	}

	return appUserID, nil
}
