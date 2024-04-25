package main

import (
	"flag"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/rs/zerolog/log"
	"main/api"
	"main/config"
	"main/services/accounts"
	"main/services/auth"
	"main/services/balances"
	"main/services/entity"
	"main/services/net_worth"
	"main/services/sync"
	"main/services/transaction"
	"main/store"
	"net/http"
	"os"
)

func main() {
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			Region: aws.String("eu-west-1"),
		},
	}))

	environmentPtr := flag.String("environment", "", "Specify the environment local, prod, docker-local)")
	flag.Parse()
	environment := *environmentPtr

	if environment != "local" && environment != "prod" && environment != "docker-local" {
		log.Error().Msg("Invalid environment specified. Please use 'local' or 'prod' or 'docker-local' as the environment.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Load the configuration with the AWS session
	_, err := config.LoadConfig(environment, awsSession)
	if err != nil {
		log.Error().Err(err).Msg("Error loading config")
		os.Exit(1)
	}

	migrationsPath := "store/migrations/" // Path to your migrations directory

	db, err := store.NewDB()
	if err != nil {
		log.Error().Err(err).Msg("Could not create DB")
		os.Exit(1)
	}
	defer db.Close()

	// Apply migrations
	if err := store.MigrateDB(migrationsPath); err != nil {
		log.Error().Err(err).Msg("Could not migrate DB")
		os.Exit(1)
	}

	userRepository := store.NewUserRepository(db)
	refreshTokenRepository := store.NewRefreshTokenRepository(db)
	bankRepository := store.NewBankRepository(db)
	entityRepository := store.NewEntityRepository(db)
	accountRepository := store.NewAccountRepository(db)
	balanceRepository := store.NewBalanceRepository(db)
	netWorthRepository := store.NewNetWorthRepository(db)

	// Create user service
	userService := auth.NewUserService(userRepository)
	transactionService := transaction.NewTransactionService()
	refreshTokenService := auth.NewRefreshTokenService(refreshTokenRepository)
	bankService := entity.NewBankService(bankRepository)
	entityService := entity.NewEntityService(entityRepository, bankService)
	accountService := accounts.NewAccountService(accountRepository)
	balanceService := balances.NewBalanceService(balanceRepository)
	netWorthService := net_worth.NewNetWorthService(accountRepository, balanceService, netWorthRepository)
	syncService := sync.NewSyncService(accountService, entityService, balanceService)

	port := "8080"
	router := api.NewRouter(
		userService,
		transactionService,
		refreshTokenService,
		entityService,
		accountService,
		netWorthService,
		syncService,
	)

	log.Info().Msgf("Server is running on :%s", port)
	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		return
	}
}
