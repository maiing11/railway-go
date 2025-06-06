package model

import "github.com/jackc/pgx/v5/pgtype"

// OpenAPI Components:
//
// components:
//   schemas:
//     Route:
//       type: object
//       properties:
//         id:
//           type: integer
//           format: int64
//         source_station:
//           type: string
//         destination_station:
//           type: string
//         travel_time:
//           type: integer
//           format: int32
//         created_at:
//           type: string
//           format: date-time
//         updated_at:
//           type: string
//           format: date-time
//       required:
//         - id
//         - source_station
//         - destination_station
//         - travel_time
//         - created_at
//         - updated_at
//     RouteRequest:
//       type: object
//       properties:
//         source_station:
//           type: string
//           maxLength: 4
//         destination_station:
//           type: string
//           maxLength: 4
//         travel_time:
//           type: integer
//           format: int32
//       required:
//         - source_station
//         - destination_station
type Route struct {
	ID                 int64            `json:"id"`
	SourceStation      string           ` json:"source_station"`
	DestinationStation string           `json:"destination_station"`
	TravelTime         int32            `json:"travel_time"`
	CreatedAt          pgtype.Timestamp `json:"created_at"`
	UpdatedAt          pgtype.Timestamp `json:"updated_at"`
}

type RouteRequest struct {
	SourceStation      string `json:"source_station" validate:"required,max=4"`
	DestinationStation string `json:"destination_station" validate:"required,max=4"`
	TravelTime         int32  `json:"travel_time"`
}
