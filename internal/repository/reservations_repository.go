package repository

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type ReservationRepository interface {
	LockSeat(ctx context.Context, scheduleID, wagonID, seatID int64, duration time.Duration) error
	UnlockSeat(ctx context.Context, scheduleID, wagonID, seatID int64) error
}

const seatLock = "seat_lock:%d:%d:%d"

func (r *redisRepository) LockSeat(ctx context.Context, scheduleID, wagonID, seatID int64, duration time.Duration) error {
	holdKey := fmt.Sprintf(seatLock, scheduleID, wagonID, seatID)
	exists, err := r.RedisClient.Exists(ctx, holdKey).Result()
	if err != nil {
		return err
	}

	if exists > 0 {
		return errors.New("seat already locked")
	}

	return r.RedisClient.Set(ctx, holdKey, "locked", duration).Err()
}

func (r *redisRepository) UnlockSeat(ctx context.Context, scheduleID, wagonID, seatID int64) error {
	holdKey := fmt.Sprintf(seatLock, scheduleID, wagonID, seatID)
	return r.RedisClient.Del(ctx, holdKey).Err()
}
