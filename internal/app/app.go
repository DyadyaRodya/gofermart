package app

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/DyadyaRodya/gofermart/internal/config"
	domainmodels "github.com/DyadyaRodya/gofermart/internal/domain/models"
	domainservices "github.com/DyadyaRodya/gofermart/internal/domain/services"
	accrualgateway "github.com/DyadyaRodya/gofermart/internal/gateways/accrual"
	"github.com/DyadyaRodya/gofermart/internal/http/controllers"
	"github.com/DyadyaRodya/gofermart/internal/http/middlewares"
	"github.com/DyadyaRodya/gofermart/internal/interactors"
	pgxrepo "github.com/DyadyaRodya/gofermart/internal/repositories/pgx"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type App struct {
	appConfig            *config.Config
	r                    *chi.Mux
	appLogger            *zap.Logger
	appStorage           *pgxrepo.StorePGX
	orderProcessorRunner func()
	close                func()

	srv *http.Server
	wg  *sync.WaitGroup
}

func InitApp(defaultServerAddress, defaultAccrualServerAddress, defaultLogLevel string, defaultSaltSize int) *App {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	appConfig := config.InitConfigFromCMD(defaultServerAddress, defaultAccrualServerAddress, defaultLogLevel, defaultSaltSize)

	appLogger, loggerMW, err := middlewares.NewLoggerMiddleware(appConfig.LogLevel)
	if err != nil {
		log.Printf("Config %+v\n", *appConfig)
		log.Fatalf("Cannot initialize logger %+v\n", err)
	}

	appLogger.Info("Config", zap.Any("config", appConfig))

	secretKeyString := os.Getenv("SECRET_KEY")
	var secretKey []byte
	if secretKeyString == "" {
		secretKey = newSecretKey(32)

		base64Text := make([]byte, base64.URLEncoding.EncodedLen(len(secretKey)))
		base64.URLEncoding.Encode(base64Text, secretKey)
		appLogger.Debug("New secret key", zap.ByteString("SECRET_KEY", base64Text))
	} else {
		appLogger.Debug("old secret key", zap.String("SECRET_KEY", secretKeyString))

		secretKey = make([]byte, base64.URLEncoding.DecodedLen(len(secretKeyString))-1)
		n, err := base64.URLEncoding.Decode(secretKey, []byte(secretKeyString))
		appLogger.Debug("after decoding secret key", zap.Int("n", n), zap.Error(err))
	}

	r := chi.NewRouter()

	// Middleware
	r.Use(loggerMW)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))

	// domain services
	jwtService := domainservices.NewJWTDomainService(secretKey)
	loginService := domainservices.NewLoginDomainService(appConfig.LoginConfig)
	luhnService := domainservices.NewLuhnDomainService()
	passwordService := domainservices.NewPasswordDomainService(appConfig.PasswordComplexityConfig)
	uuidGeneratorService := domainservices.NewUUID4Generator()

	// gateway
	accrualGateway, err := accrualgateway.NewGateway(appConfig.AccrualGatewayConfig, func(msg string) {
		appLogger.Error(msg) // pass error logger function to log some errors
	})
	if err != nil {
		appLogger.Fatal("Cannot initialize accrual gateway", zap.Error(err))
	}

	// repository
	pool, err := pgxpool.New(ctx, appConfig.DSN)
	if err != nil {
		appLogger.Fatal("Cannot create connection pool to database", zap.String("DATABASE_URI", appConfig.DSN), zap.Error(err))
	}
	closeStorage := pool.Close
	repo := pgxrepo.NewStorePGX(pool, appLogger)

	// auth middleware
	authMW := middlewares.NewAuthMiddleware(jwtService, repo)

	// chan
	processOrderChan := make(chan *domainmodels.Order, 1024)

	// interactors
	loginInteractor := interactors.NewLoginInteractor(repo, passwordService, appConfig.SaltSize)
	registerInteractor := interactors.NewRegisterInteractor(repo, passwordService, loginService, uuidGeneratorService, appConfig.SaltSize)

	addOrderInteractor := interactors.NewAddOrderInteractor(repo, luhnService, processOrderChan)
	getOrdersInteractor := interactors.NewGetOrdersInteractor(repo)
	processOrderInteractor := interactors.NewProcessOrderInteractor(repo, accrualGateway)

	getWithdrawalsInteractor := interactors.NewGetWithdrawalsInteractor(repo)
	withdrawInteractor := interactors.NewWithdrawInteractor(repo, luhnService)

	getBalanceInteractor := interactors.NewGetBalanceInteractor(repo)

	// controllers
	loginController := controllers.NewLoginController(jwtService, loginInteractor)
	registerController := controllers.NewRegisterController(jwtService, registerInteractor)

	addOrderController := controllers.NewAddOrderController(addOrderInteractor)
	getOrdersController := controllers.NewGetUserOrdersController(getOrdersInteractor)

	getWithdrawalsController := controllers.NewGetUserWithdrawalsController(getWithdrawalsInteractor)
	withdrawController := controllers.NewWithdrawController(withdrawInteractor)

	getBalanceController := controllers.NewGetBalanceController(getBalanceInteractor)

	// routing
	initRoutes(
		r, authMW,
		loginController,
		registerController,
		addOrderController,
		getOrdersController,
		getWithdrawalsController,
		withdrawController,
		getBalanceController,
	)

	// separate processor into other routine
	ctxProcessor, cancelProcessor := context.WithCancel(ctx)
	processorRunner := func() {
		appLogger.Info("Starting order processor")
		go processOrderInteractor.Run(ctxProcessor, func(msg string) {
			appLogger.Error(msg) // pass error logger function to log some errors
		}, processOrderChan)

		dbSess, err := repo.NewSession(ctxProcessor)
		if err != nil {
			appLogger.Fatal("Cannot initialize session for order processor", zap.Error(err))
		}
		defer dbSess.Close(ctxProcessor)
		sess := dbSess.(*pgxrepo.SessionPGX)
		orders, err := sess.GetUnprocessedOrders(ctxProcessor)
		if err != nil {
			appLogger.Fatal("Cannot get unprocessed orders", zap.Error(err))
		}
		for _, order := range orders {
			order := order
			go func() { processOrderChan <- order }()
		}
	}

	closer := func() {
		closeStorage()
		cancelProcessor()
		for range processOrderChan {

		}
		close(processOrderChan)
		cancel()
	}
	return &App{
		appConfig:            appConfig,
		r:                    r,
		appLogger:            appLogger,
		appStorage:           repo,
		close:                closer,
		orderProcessorRunner: processorRunner,
	}
}

func initRoutes(
	r *chi.Mux,
	authMW *middlewares.AuthMiddleware,
	loginController *controllers.LoginController,
	registerController *controllers.RegisterController,
	addOrderController *controllers.AddOrderController,
	getOrdersController *controllers.GetUserOrdersController,
	getWithdrawalsController *controllers.GetUserWithdrawalsController,
	withdrawController *controllers.WithdrawController,
	getBalanceController *controllers.GetBalanceController,
) {
	r.Route("/api", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Post("/login", loginController.ServeHTTP)
			r.Post("/register", registerController.ServeHTTP)

			r.Route("/orders", func(r chi.Router) {
				r.Post("/", authMW.WithAuth(addOrderController))
				r.Get("/", authMW.WithAuth(getOrdersController))
			})

			r.Route("/balance", func(r chi.Router) {
				r.Get("/", authMW.WithAuth(getBalanceController))
				r.Post("/withdraw", authMW.WithAuth(withdrawController))
			})

			r.Get("/withdrawals", authMW.WithAuth(getWithdrawalsController))
		})
	})
}

func newSecretKey(size int) []byte {
	secretKey := make([]byte, size)
	_, err := rand.Read(secretKey)
	if err != nil {
		panic(err)
	}
	return secretKey
}

func (a *App) Run() error {
	a.wg = &sync.WaitGroup{}

	a.wg.Add(1)

	a.appLogger.Info("Initializing storage")
	ctx := context.Background()
	ctx, cancelFunc := context.WithTimeout(ctx, 3*time.Second)
	defer cancelFunc()
	err := a.appStorage.InitSchema(ctx)
	if err != nil {
		a.appLogger.Fatal("appStorage.InitSchema error", zap.Error(err))
		return err
	}

	a.orderProcessorRunner()

	a.srv = &http.Server{Addr: a.appConfig.ServerAddress, Handler: a.r}
	a.appLogger.Info("Starting server at", zap.String("address", a.appConfig.ServerAddress))

	defer a.wg.Done() // let know we are done cleaning up
	if err = a.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		a.appLogger.Fatal("ListenAndServe error", zap.Error(err))
		return err
	}
	return nil
}

func (a *App) Shutdown(signal os.Signal) error {
	defer a.close()

	ctx := context.Background()
	a.appLogger.Info("Stopped server on signal", zap.String("signal", signal.String()))

	ctx, cancelFunc := context.WithTimeout(ctx, 10*time.Second)
	defer cancelFunc()
	err := a.srv.Shutdown(ctx)
	if err != nil {
		a.appLogger.Fatal("Server Shutdown error", zap.Error(err))
	}

	a.wg.Wait()
	return nil
}
