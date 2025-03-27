package repository

import (
	"context"
	"fmt"
	"time"
)

type ReservationRepository interface {
	LockSeat(ctx context.Context, scheduleID, wagonID, seatID int64, duration time.Duration) (bool, error)
	UnlockSeat(ctx context.Context, scheduleID, wagonID, seatID int64) error
}

func (r *redisRepository) LockSeat(ctx context.Context, scheduleID, wagonID, seatID int64, duration time.Duration) (bool, error) {
	holdKey := fmt.Sprintf("seat_lock:%d:%d:%d", scheduleID, wagonID, seatID)
	success, err := r.RedisClient.SetNX(ctx, holdKey, "locked", duration).Result()
	if err != nil {
		return false, err
	}

	return success, nil
}

func (r *redisRepository) UnlockSeat(ctx context.Context, scheduleID, wagonID, seatID int64) error {
	holdKey := fmt.Sprintf("seat_lock:%d:%d:%d", scheduleID, wagonID, seatID)
	return r.RedisClient.Del(ctx, holdKey).Err()
}
