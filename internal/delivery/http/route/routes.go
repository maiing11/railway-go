package route

import (
	"railway-go/internal/delivery/http"
	"railway-go/internal/delivery/http/middleware"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                   *fiber.App
	UserController        http.UserControllers
	ReservationController http.ReservationControllers
	PassengerController   http.PassengerControllers
	ScheduleController    http.ScheduleControllers
	PaymentController     http.PaymentControllers
	RouteController       http.RouteControllers
	SeatController        http.SeatControllers
	TrainController       http.TrainControllers
	WagonController       http.WagonControllers
	DiscountController    http.DiscountControllers
	StationController     http.StationControllers
	AuthMiddleware        *middleware.AuthMiddleware
}

func (c *RouteConfig) Setup() {
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupAuthRoute() {
	// Public routes
	c.App.Post("/users/register", c.UserController.Register)
	c.App.Post("/users/register/_admin", c.UserController.RegisterAdmin)
	c.App.Post("/users/login", c.UserController.Login)
	c.App.Post("/users/logout", c.UserController.Logout)

	// Authenticated user routes
	auth := c.App.Group("/auth", c.AuthMiddleware.AuthRequired())
	auth.Post("/reservations", c.ReservationController.CreateReservation)
	auth.Get("/reservations", c.ReservationController.GetDetailReservation)
	auth.Delete("/reservations", c.ReservationController.DeleteReservation)
	auth.Put("/reservations/_canceled", c.ReservationController.CancelReservation)
	auth.Post("/reservations/payments", c.PaymentController.MockPaymentWebhook)

	auth.Get("/schedules", c.ScheduleController.GetSchedule)
	auth.Get("/schedules/search", c.ScheduleController.SearchSchedules)

	auth.Post("/passengers", c.PassengerController.CreatePassenger)
	auth.Get("/passengers", c.PassengerController.GetPassenger)
	auth.Put("/passengers", c.PassengerController.UpdatePassenger)
	auth.Delete("/passengers", c.PassengerController.DeletePassenger)

	// Admin routes
	admin := c.App.Group("/admin", c.AuthMiddleware.AuthRequired(), c.AuthMiddleware.AdminOnly())
	admin.Get("/reservations", c.ReservationController.GetAllReservations)

	// General Affairs routes
	ga := c.App.Group("/ga", c.AuthMiddleware.AuthRequired(), c.AuthMiddleware.GeneralAffairs())
	ga.Post("/schedules", c.ScheduleController.CreateSchedule)
	ga.Put("/schedules", c.ScheduleController.UpdateSchedule)
	ga.Delete("/schedules", c.ScheduleController.DeleteSchedule)

	ga.Post("/train_routes", c.RouteController.CreateRoute)
	ga.Get("/train_routes", c.RouteController.GetRoute)
	ga.Put("/train_routes", c.RouteController.UpdateRoute)
	ga.Delete("/train_routes", c.RouteController.DeleteRoute)
	ga.Get("/train_routes/list", c.RouteController.GetRoutes)

	ga.Post("/train_seats", c.SeatController.CreateSeat)
	ga.Get("/train_seats", c.SeatController.GetSeat)
	ga.Get("/train_seats", c.SeatController.GetSeats)
	ga.Put("/train_seats", c.SeatController.UpdateSeat)
	ga.Delete("/train_seats/list", c.SeatController.DeleteSeat)

	ga.Post("/trains", c.TrainController.CreateTrain)
	ga.Get("/trains", c.TrainController.GetTrain)
	ga.Get("/trains", c.TrainController.GetTrains)
	ga.Put("/trains", c.TrainController.UpdateTrain)
	ga.Delete("/trains/list", c.TrainController.DeleteTrain)

	ga.Post("/train_wagons", c.WagonController.CreateWagon)
	ga.Get("/train_wagons", c.WagonController.GetWagon)
	ga.Get("/train_wagons/list", c.WagonController.GetWagons)
	ga.Put("/train_wagons", c.WagonController.UpdateWagon)
	ga.Delete("/train_wagons", c.WagonController.DeleteWagon)

	ga.Post("/train_discounts", c.DiscountController.CreateDiscount)
	ga.Get("/train_discounts", c.DiscountController.GetDiscount)

	ga.Post("/train_stations", c.StationController.CreateStation)
	ga.Put("/train_stations/set", c.StationController.UpdateStation)
	ga.Get("/train_stations/search", c.StationController.GetStationByName)
	ga.Get("/train_stations/list", c.StationController.GetStations)
	ga.Get("/train_stations/:code", c.StationController.GetStationByCode)
	ga.Delete("/train_stations", c.StationController.DeleteStation)
}
