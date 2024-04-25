package api

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"main/models"
	"main/services/auth"
	"net/http"
)

func registerHandler(userService *auth.UserService, refreshService *auth.RefreshTokenService) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// Parse the request body to get the user data.
		var user *models.UserRequest

		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Use the validator library to validate the UserRequest.
		validate := validator.New()
		if err := validate.Struct(user); err != nil {
			// Handle validation errors
			http.Error(w, "Validation failed", http.StatusBadRequest)
			return
		}

		// Call the user service to register the user.
		if err := userService.Register(user); err != nil {
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			return
		}

		// Generate an access token.
		authToken, err := auth.GenerateAuthToken(user.Username)
		if err != nil {
			http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
			return
		}

		// Generate a refresh token.
		refreshToken, err := refreshService.GenerateNewRefreshToken(user.Username)
		if err != nil {
			http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
			return
		}

		// Prepare the response JSON with both tokens.
		response := struct {
			AccessToken    string `json:"access_token"`
			RefreshToken   string `json:"refresh_token"`
			LeanCustomerId string `json:"lean_customer_id"`
		}{
			AccessToken:    authToken,
			RefreshToken:   refreshToken.Token,
			LeanCustomerId: user.LeanCustomerId,
		}

		// Encode the response as JSON.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

func loginHandler(userService *auth.UserService, refreshService *auth.RefreshTokenService) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// Parse the request body to get the username and password.
		var requestBody struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		authToken, user, err := userService.Login(requestBody.Email, requestBody.Password)
		if err != nil {
			http.Error(w, "Invalid login", http.StatusUnauthorized)
			return
		}

		refreshToken, err := refreshService.GetRefreshToken(requestBody.Email)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		response := struct {
			AccessToken    string `json:"access_token"`
			RefreshToken   string `json:"refresh_token"`
			LeanCustomerId string `json:"lean_customer_id"`
			IsConnected    bool   `json:"is_connected"`
		}{
			AccessToken:    authToken,
			RefreshToken:   refreshToken.Token,
			LeanCustomerId: user.LeanCustomerId,
			IsConnected:    user.ConnectedState,
		}

		// Encode the response as JSON.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
		return
	}
}
