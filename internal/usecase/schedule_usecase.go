package usecase

import (
	"context"
	"railway-go/internal/constant/model"
	"railway-go/internal/repository"
	"railway-go/internal/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type ScheduleUC interface {
	CreateSchedule(ctx context.Context, request *model.ScheduleRequest) (repository.Schedule, error)
	UpdateSchedule(ctx context.Context, reqId int64, request model.ScheduleRequest) error
	GetSchedule(ctx context.Context, id int64) (repository.Schedule, error)
	DeleteSchedule(ctx context.Context, id int64) error
	SearchSchedules(ctx context.Context, request *model.SearchScheduleRequest) ([]model.SearchScheduleResponse, error)
}

type ScheduleUsecase struct {
	*UseCase
	TrainUC
	RouteUC
}

func NewScheduleUsecase(useCase *UseCase) ScheduleUC {
	return &ScheduleUsecase{UseCase: useCase}
}

func (uc *ScheduleUsecase) CreateSchedule(ctx context.Context, request *model.ScheduleRequest) (repository.Schedule, error) {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		uc.Log.Warn("failed to begin transaction", zap.Error(err))
		return repository.Schedule{}, fiber.NewError(fiber.StatusInternalServerError, "failed to begin transaction")
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	train, err := tx.GetTrain(ctx, request.TrainID)
	if err != nil {
		uc.Log.Warn("invalid train id or train not found", zap.Error(err))
		return repository.Schedule{}, fiber.ErrInternalServerError
	}

	route, err := tx.GetRoute(ctx, request.RouteID)
	if err != nil {
		uc.Log.Warn("invalid route id or route not found ", zap.Error(err))
		return repository.Schedule{}, fiber.ErrInternalServerError
	}

	// validate request
	if err := uc.Validate.Struct(request); err != nil {
		uc.Log.Warn("error validating request body", zap.Error(err))
		return repository.Schedule{}, fiber.ErrBadRequest
	}

	arrival := &pgtype.Timestamp{
		Time:  request.DepartureDate.Time.Add(time.Duration(int64(route.TravelTime)) * time.Minute),
		Valid: true,
	}

	schedule := repository.CreateScheduleParams{
		TrainID:        train.ID,
		DepartureDate:  request.DepartureDate,
		ArrivalDate:    *arrival,
		AvailableSeats: request.AvailableSeats,
		Price:          request.Price,
		RouteID:        route.ID,
	}

	response, err := tx.CreateSchedule(ctx, schedule)
	if err != nil {
		uc.Log.Warn("Error creating schedule", zap.Error(err))
		return repository.Schedule{}, fiber.NewError(fiber.StatusInternalServerError, "error creating schedule")
	}

	// commit transaction
	if err := tx.Commit(ctx); err != nil {
		uc.Log.Error("failed to commit transaction", zap.Error(err))
		return repository.Schedule{}, fiber.ErrInternalServerError
	}

	return response, nil
}

func (uc *ScheduleUsecase) UpdateSchedule(ctx context.Context, reqId int64, request model.ScheduleRequest) error {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		uc.Log.Warn("failed to begin transaction", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "failed to begin transaction")
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	s, err := tx.GetSchedule(ctx, reqId)
	if err != nil {
		uc.Log.Warn("error getting schedule", zap.Error(err))
		return fiber.ErrNotFound
	}

	train, err := tx.GetTrain(ctx, request.TrainID)
	if err != nil {
		uc.Log.Warn("invalid train id or train not found", zap.Error(err))
		return fiber.ErrInternalServerError
	}

	route, err := tx.GetRoute(ctx, request.RouteID)
	if err != nil {
		uc.Log.Warn("invalid route id or route not found ", zap.Error(err))
		return fiber.ErrInternalServerError
	}

	// validate request
	if err := uc.Validate.Struct(request); err != nil {
		uc.Log.Warn("error validating request body")
		return fiber.ErrBadRequest
	}

	arrival := &pgtype.Timestamp{
		Time:  request.DepartureDate.Time.Add(time.Duration(int64(route.TravelTime)) * time.Minute),
		Valid: true,
	}

	schedule := repository.UpdateScheduleParams{
		ID:             s.ID,
		TrainID:        train.ID,
		RouteID:        route.ID,
		DepartureDate:  request.DepartureDate,
		ArrivalDate:    *arrival,
		Price:          request.Price,
		AvailableSeats: request.AvailableSeats,
	}

	if err := tx.UpdateSchedule(ctx, schedule); err != nil {
		uc.Log.Warn("error updating schedule", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "error updating schedule")
	}

	if err := tx.Commit(ctx); err != nil {
		uc.Log.Error("error commiting transaction", zap.Error(err))
		return fiber.ErrInternalServerError
	}

	return nil
}

func (uc *ScheduleUsecase) GetSchedule(ctx context.Context, id int64) (repository.Schedule, error) {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		uc.Log.Warn("failed to begin transaction", zap.Error(err))
		return repository.Schedule{}, fiber.NewError(fiber.StatusInternalServerError, "failed to begin transaction")
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	schedule, err := tx.GetSchedule(ctx, id)
	if err != nil {
		uc.Log.Warn("failed to find schedule", zap.Error(err))
		return repository.Schedule{}, fiber.NewError(fiber.StatusInternalServerError, "failed to find schedule")
	}

	if err := tx.Commit(ctx); err != nil {
		uc.Log.Error("failed to commit transaction", zap.Error(err))
		return repository.Schedule{}, fiber.ErrInternalServerError
	}

	return schedule, nil
}

func (uc *ScheduleUsecase) DeleteSchedule(ctx context.Context, id int64) error {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		uc.Log.Warn("failed to begin transaction", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "failed to begin transaction")
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	schedule, err := tx.GetSchedule(ctx, id)
	if err != nil {
		uc.Log.Warn("error schedule not found", zap.Error(err))
		return fiber.NewError(fiber.StatusNotFound, "error schedule not found")
	}

	if err := tx.DeleteSchedule(ctx, schedule.ID); err != nil {
		uc.Log.Error("failed to delete schedule", zap.Error(err))
		return nil
	}

	if err := tx.Commit(ctx); err != nil {
		uc.Log.Error("failed to commit transaction", zap.Error(err))
		return fiber.ErrInternalServerError
	}

	return nil
}

func (uc *ScheduleUsecase) SearchSchedules(ctx context.Context, request *model.SearchScheduleRequest) ([]model.SearchScheduleResponse, error) {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		uc.Log.Warn("failed to begin transaction", zap.Error(err))
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to begin transaction")
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	parsedTime, err := time.Parse("2006-01-02", request.DepartureDate)
	if err != nil {
		return nil, utils.WrapError(fiber.StatusBadRequest, uc.Log, utils.Warn, err, "failed to convert departuredate")
	}

	// check if departure date is in the past
	if parsedTime.Before(time.Now().Truncate(24 * time.Hour)) {
		return nil, fiber.NewError(fiber.StatusBadRequest, "departure date cannot be in the past")
	}

	response, err := tx.SearchSchedules(ctx, repository.SearchSchedulesParams{
		Column1:       request.SourceStation,
		Column2:       request.DestinationStation,
		DepartureDate: pgtype.Timestamp{Time: parsedTime, Valid: true},
	})
	if err != nil {
		uc.Log.Warn("failed to search schedule", zap.Error(err))
		return nil, fiber.NewError(fiber.StatusNotFound, "failed to search schedule")
	}

	if err := tx.Commit(ctx); err != nil {
		uc.Log.Error("faield to commit transaction", zap.Error(err))
		return nil, fiber.ErrInternalServerError
	}

	// map response to model
	var schedules []model.SearchScheduleResponse
	for _, schedule := range response {
		// return only future schedules
		if schedule.DepartureDate.Time.After(time.Now()) {
			schedules = append(schedules, model.SearchScheduleResponse{
				ScheduleID:         schedule.ScheduleID,
				TrainName:          schedule.TrainName,
				SourceStation:      schedule.SourceStation,
				DestinationStation: schedule.DestinationStation,
				DepartureDate:      schedule.DepartureDate,
				ArrivalDate:        schedule.ArrivalDate,
				AvailableSeats:     schedule.AvailableSeats,
				Price:              schedule.Price,
			})
		}
	}

	return schedules, nil
}
