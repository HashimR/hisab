package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"main/config"
	"main/models"
	"main/store"
	"main/store/dbmodels"
	"time"
)

type UserService struct {
	userRepository *store.UserRepository
}

func NewUserService(ur *store.UserRepository) *UserService {
	return &UserService{
		userRepository: ur,
	}
}

func (us *UserService) Register(user *models.UserRequest) error {
	// Hash the password before storing it.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	openBankingId := generateOpenBankingID(user.Username)
	leanCustomerId, err := CreateLeanCustomer(openBankingId)
	if err != nil {
		return fmt.Errorf("could not connect to Lean to register user")
	}

	// Create a UserData struct with the provided information.
	userData := &dbmodels.UserData{
		Email:          user.Username,
		Password:       string(hashedPassword),
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		PhoneNumber:    user.PhoneNumber,
		Country:        user.Country,
		OpenBankingId:  openBankingId,
		LeanCustomerId: leanCustomerId,
	}

	// Call the UserRepository to create the user.
	if err := us.userRepository.CreateUser(userData); err != nil {
		return err
	}

	user.LeanCustomerId = leanCustomerId
	return nil
}

func generateOpenBankingID(email string) string {
	// Generate a unique hash using timestamp and email.
	timestamp := time.Now().UnixNano()
	uniqueIdentifier := fmt.Sprintf("%s:%d", email, timestamp)

	// Create a SHA-256 hash of the concatenated data.
	hasher := sha256.New()
	hasher.Write([]byte(uniqueIdentifier))
	hash := hex.EncodeToString(hasher.Sum(nil))

	// TODO: check if already exists

	return hash
}

func (us *UserService) Login(email, password string) (string, *dbmodels.UserData, error) {
	// Retrieve the user by email from the UserRepository.
	user, err := us.userRepository.GetUserByEmail(email)
	if err != nil {
		return "", nil, err // UserRequest not found or error occurred.
	}

	// Compare the provided password with the stored hashed password.
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", nil, err
	}
	authToken, err := GenerateAuthToken(user.Email)
	return authToken, user, err
}

func (us *UserService) GetUserByEmail(email string) (*dbmodels.UserData, error) {
	// Retrieve the user by email from the UserRepository.
	user, err := us.userRepository.GetUserByEmail(email)
	if err != nil {
		return nil, err // UserRequest not found or error occurred.
	}
	return user, nil
}

func (us *UserService) GetUserByLeanCustomerId(leanCustomerId string) (*models.User, error) {
	userData, err := us.userRepository.GetUserByLeanCustomerId(leanCustomerId)
	if err != nil {
		log.Error().Err(err).Msgf("Could not get user by lean customer ID: %s", leanCustomerId)
		return nil, err
	}

	return &models.User{
		ID:             userData.ID,
		Email:          userData.Email,
		FirstName:      userData.FirstName,
		LastName:       userData.LastName,
		PhoneNumber:    userData.PhoneNumber,
		Country:        userData.Country,
		OpenBankingId:  userData.OpenBankingId,
		LeanCustomerId: userData.LeanCustomerId,
		ConnectedState: userData.ConnectedState,
	}, nil
}

func (us *UserService) StoreConnectedState(leanCustomerId string) error {
	return us.userRepository.StoreConnectedState(leanCustomerId)
}

func GenerateAuthToken(userID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix() // Token expiration time

	jwtKey := config.GetConfig().GetString("jwt-key")
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateAuthToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method and return the secret key.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		jwtKey := config.GetConfig().GetString("jwt-key")
		return []byte(jwtKey), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["user_id"].(string)
		return userID, nil
	} else {
		return "", fmt.Errorf("invalid token")
	}
}
