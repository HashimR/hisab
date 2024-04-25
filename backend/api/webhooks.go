package api

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
	"main/models"
	"main/services/accounts"
	"main/services/auth"
	"main/services/entity"
	"net/http"
)

const EntityCreated = "entity.created"

func getEntityHandler(es *entity.EntityService, us *auth.UserService, as *accounts.Service) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// Parse the request body to get the user data.
		var entityModel models.EntityResponse

		if err := json.NewDecoder(r.Body).Decode(&entityModel); err != nil {
			log.Error().Err(err).Msg("Invalid request body")
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if entityModel.Type != EntityCreated {
			log.Error().Msgf("Invalid request type: %s", entityModel.Type)
			http.Error(w, "Invalid request type", http.StatusBadRequest)
			return
		}

		validate := validator.New()
		if err := validate.Struct(entityModel); err != nil {
			// Handle validation errors
			http.Error(w, "Validation failed", http.StatusBadRequest)
			return
		}

		user, err := us.GetUserByLeanCustomerId(entityModel.Payload.LeanCustomerId)
		if err != nil {
			log.Error().Err(err).Msgf("Could not get user for lean customer id: %s", entityModel.Payload.LeanCustomerId)
			return
		}

		if err := es.StoreEntity(entityModel, user.ID); err != nil {
			log.Error().Err(err).Msgf("Failed to store entityModel: %s", entityModel)
			http.Error(w, "Failed to store entityModel", http.StatusInternalServerError)
			return
		}

		_, err = as.FetchAndStoreAccounts(entityModel.Payload.Id, user.ID)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to store lean accounts")
			http.Error(w, "Failed to store lean accounts", http.StatusInternalServerError)
			return
		}

		//TODO: find a way to get balances & transactions (and how often to retrieve them)

		log.Info().Msgf("Storing connected state for: %s", entityModel.Payload.LeanCustomerId)
		_ = us.StoreConnectedState(entityModel.Payload.LeanCustomerId)

		// Encode the response as JSON.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
	}
}
