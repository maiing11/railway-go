-- name: GetDiscountByCode :one
SELECT * FROM discount_codes
WHERE code = $1 AND expires_at > NOW() AND max_uses > 0;

-- name: ApplyDiscountToReservation :exec
INSERT INTO reservation_discounts (reservation_id, discount_id)
VALUES ($1, $2);

-- name: ReduceDiscountUsage :exec
UPDATE discount_codes
SET max_uses = max_uses - 1
WHERE id = $1;

-- name: GetDiscountsForReservation :many
SELECT dc.id, dc.code, dc.discount_percent, dc.expires_at, dc.max_uses
FROM discount_codes dc
JOIN reservation_discounts rd ON dc.id = rd.discount_id
WHERE rd.reservation_id = $1;

-- name: CreateDiscountCode :one
INSERT INTO discount_codes (code, discount_percent, expires_at, max_uses)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateDiscountCode :exec
UPDATE discount_codes
SET code = $2,
discount_percent = $3,
expires_at = $4,
max_uses = $5,
updated_at = NOW()
WHERE id = $1;