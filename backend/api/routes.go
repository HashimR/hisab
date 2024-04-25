package api

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
	"main/middleware"
	"main/services/accounts"
	"main/services/auth"
	"main/services/entity"
	"main/services/net_worth"
	"main/services/sync"
	"main/services/transaction"
	"net/http"
	"strings"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Implement your own origin check logic here, for example:
		// return r.Header.Get("Origin") == "http://yourdomain.com"
		return true // Allow connections from any origin
	},
	// Add other upgrader configurations if necessary
}

func NewRouter(
	userService *auth.UserService,
	transactionService *transaction.Service,
	refreshTokenService *auth.RefreshTokenService,
	entityService *entity.EntityService,
	accountService *accounts.Service,
	netWorthService *net_worth.Service,
	syncService *sync.Service,
) http.Handler {
	router := httprouter.New()

	// User registration and login routes
	router.POST("/register", registerHandler(userService, refreshTokenService))
	router.POST("/login", loginHandler(userService, refreshTokenService))
	router.POST("/refresh", getRefreshHandler(refreshTokenService))
	router.POST("/entity", getEntityHandler(entityService, userService, accountService))
	router.GET("/transactions", middleware.AuthenticationMiddleware(getTransactionsHandler(transactionService), userService))
	router.GET("/net-worth", middleware.AuthenticationMiddleware(getNetWorthHandler(netWorthService), userService))
	router.POST("/net-worth", middleware.AuthenticationMiddleware(getNetWorthCalculateHandler(netWorthService), userService))
	router.GET("/accounts", middleware.AuthenticationMiddleware(getAccountsHandler(accountService), userService))

	router.GET("/sync", middleware.AuthenticationMiddleware(getSyncHandler(syncService), userService))

	//router.GET("/transactions", getTransactionsHandler(transactionService))

	return router
}

func getSyncHandler(syncService *sync.Service) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// Upgrade HTTP connection to WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error().Err(err).Msg("Error upgrading to WebSocket")
			http.Error(w, "Error upgrading to WebSocket", http.StatusInternalServerError)
			return
		}
		defer func() {
			if err := conn.Close(); err != nil {
				log.Error().Err(err).Msg("Error closing WebSocket connection")
			}
		}()

		userID, ok := r.Context().Value("userId").(int)
		if !ok {
			log.Error().Msg("Error retrieving userID from context")
			http.Error(w, "Internal error. Please try again.", http.StatusInternalServerError)
			return
		}

		//// Read message from client
		//_, msg, err := conn.ReadMessage()
		//if err != nil {
		//	log.Error().Err(err).Msg("Error reading message")
		//	http.Error(w, "Internal error. Please try again.", http.StatusInternalServerError)
		//	return
		//}
		//log.Info().Int("userID", userID).Str("message", string(msg)).Msg("Received message from user")

		err = syncService.Sync(userID)
		if err != nil {
			log.Error().Err(err).Msg("Could not sync")
			http.Error(w, "Internal error. Please try again.", http.StatusInternalServerError)
		}

		// for all entities, fetch accounts, balances, transactions. Update DB

		// Close WebSocket connection after sync is complete or in case of error
		if err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
			log.Error().Err(err).Msg("Error sending close message")
		}
	}
}

func getAccountsHandler(accountService *accounts.Service) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		userId := r.Context().Value("userID").(int)

		netWorthGraph, err := accountService.GetAccounts(userId)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(netWorthGraph)
	}
}

func getNetWorthCalculateHandler(netWorthService *net_worth.Service) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		userId := r.Context().Value("userID").(int)

		netWorthGraph, err := netWorthService.CalculateCurrentNetWorth(userId)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(netWorthGraph)
	}
}

func getNetWorthHandler(netWorthService *net_worth.Service) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		userId := r.Context().Value("userID").(int)

		netWorthGraph, err := netWorthService.GetLastXNetWorthRecords(userId, 5)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(netWorthGraph)
	}
}

func getRefreshHandler(refreshTokenService *auth.RefreshTokenService) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// Get the refresh token from the Authorization header.
		authHeader := r.Header.Get("Authorization")

		// Check if the Authorization header is missing or doesn't have the expected format.
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Extract the refresh token from the header (assuming it's in the format "Bearer <token>").
		refreshTokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the refresh token.
		refreshToken, err := refreshTokenService.ValidateRefreshToken(refreshTokenString)
		if err != nil {
			http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
			return
		}

		// Generate a new JWT for the user.
		authToken, err := auth.GenerateAuthToken(refreshToken.Username)
		if err != nil {
			http.Error(w, "Failed to generate new JWT", http.StatusInternalServerError)
			return
		}

		// Respond with the new JWT.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := struct {
			Token        string `json:"auth-token"`
			RefreshToken string `json:"refresh-token"`
		}{
			Token:        authToken,
			RefreshToken: refreshToken.Token,
		}

		if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

func getTransactionsHandler(transactionService *transaction.Service) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		transactions, err := transactionService.GetTransactions()

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(transactions)
	}
}
