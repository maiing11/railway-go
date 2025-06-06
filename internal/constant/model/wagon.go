package model

type WagonRequest struct {
	TrainID     int64  `json:"train_id" validate:"required"`
	WagonNumber int32  `json:"wagon_number" validate:"required"`
	ClassType   string `json:"class_type" validate:"required"`
	TotalSeats  int32  `json:"total_seats" validate:"required"`
}
