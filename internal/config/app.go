package config

import (
	"context"
	"railway-go/internal/delivery/http"
	"railway-go/internal/delivery/http/middleware"
	"railway-go/internal/delivery/http/route"
	"railway-go/internal/repository"
	"railway-go/internal/usecase"
	"railway-go/internal/utils/token"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type BootstrapConfig struct {
	DB          *pgxpool.Pool
	App         *fiber.App
	Log         *zap.Logger
	Validate    *validator.Validate
	Config      *viper.Viper
	TokenMaker  token.Maker
	RedisClient *redis.Client
}

// do bootstrap here
func Boostrap(config BootstrapConfig) {
	// setup repositories
	repo := repository.NewStore(config.DB, config.RedisClient)

	baseUsecase := usecase.NewUseCase(repo, config.Log, config.Validate)

	// setup usecases
	userSessionUC := usecase.NewUserSessionUsecase(baseUsecase, config.TokenMaker, config.Config)
	reservationUC := usecase.NewReservationUsecase(baseUsecase)
	scheduleUC := usecase.NewScheduleUsecase(baseUsecase)
	paymentUC := usecase.NewPaymentUsecase(baseUsecase)
	discountUC := usecase.NewDiscountUsecase(baseUsecase)
	passengerUC := usecase.NewPassengerUsecase(baseUsecase)
	routeUC := usecase.NewRouteUsecase(baseUsecase)
	seatUC := usecase.NewSeatUsecase(baseUsecase)
	trainUC := usecase.NewTrainUsecase(baseUsecase)
	wagonUC := usecase.NewWagonUsecase(baseUsecase)
	stationUC := usecase.NewStationUsecase(baseUsecase)

	StartReservationCleanup(reservationUC, paymentUC, config.Log)
	// setup controlers
	userSesionController := http.NewUserSessionController(userSessionUC, config.Log)
	reservationController := http.NewReservationController(reservationUC, config.Log, userSessionUC)
	scheduleController := http.NewScheduleController(scheduleUC, config.Log)
	paymentController := http.NewPaymentController(paymentUC, config.Log)
	discountController := http.NewDiscountController(config.Log, discountUC)
	passengerController := http.NewPassengerController(passengerUC, userSessionUC, config.Log)
	routeController := http.NewRouteController(routeUC, config.Log)
	seatController := http.NewSeatController(config.Log, seatUC)
	trainController := http.NewTrainController(trainUC, config.Log)
	wagonController := http.NewWagonController(config.Log, wagonUC)
	stationController := http.NewStationController(stationUC, config.Log)

	// setup middlewares
	userSessionMiddlewares := middleware.NewAuthMiddleware(userSessionUC, config.TokenMaker)

	// setup routes
	routeConfig := route.RouteConfig{
		App:                   config.App,
		UserController:        userSesionController,
		ReservationController: reservationController,
		ScheduleController:    scheduleController,
		PaymentController:     paymentController,
		DiscountController:    discountController,
		PassengerController:   passengerController,
		RouteController:       routeController,
		SeatController:        seatController,
		TrainController:       trainController,
		WagonController:       wagonController,
		StationController:     stationController,
		AuthMiddleware:        userSessionMiddlewares,
	}

	routeConfig.Setup()
}

func StartReservationCleanup(reservationUC usecase.ReservationUC, paymentUC usecase.PaymentUC, log *zap.Logger) {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

			err := paymentUC.AutoCancelExpiredPayments(ctx)
			if err != nil {
				log.Error("failed to auto cancel expired payments", zap.Error(err))
			}
			err = reservationUC.AutoDeleteReservations(ctx)
			if err != nil {
				log.Error("failed to auto delete expired reservations", zap.Error(err))
			} else {
				log.Info("successfully auto deleted expired reservations")
			}
			cancel()
		}
	}()
}
