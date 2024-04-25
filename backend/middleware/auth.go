package middleware

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"main/services/auth"
	"net/http"
	"strings"
)

func AuthenticationMiddleware(next httprouter.Handle, us *auth.UserService) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Split the Authorization header into parts.
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		// Extract the token part.
		tokenString := parts[1]

		// Verify and parse the JWT token.
		email, err := auth.ValidateAuthToken(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Attach the email to the request context for later use by handlers.
		r = r.WithContext(context.WithValue(r.Context(), "email", email))

		user, err := us.GetUserByEmail(email)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), "userId", user.ID))

		// If the JWT is valid, call the next handler in the chain.
		next(w, r, ps)
	}
}
