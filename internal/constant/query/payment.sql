-- name: GetPayment :one
SELECT * FROM  payments
WHERE id = $1 LIMIT 1;

-- name: ListPayments :many
SELECT * FROM payments
ORDER BY id;

-- name: CreatePayment :exec
INSERT INTO payments (
    reservation_id, payment_method, amount, transaction_id, payment_date, gateway_response, payment_status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
);

-- name: UpdatePayment :exec
UPDATE payments
  set reservation_id = $2,
  payment_method = $3,
  amount = $4,
  transaction_id = $5,
  payment_date = $6,
  gateway_response = $7,
  payment_status = $8
WHERE id = $1;

-- name: CompletePayment :exec
UPDATE payments
  set payment_status = 'success', payment_date = NOW()
WHERE id = $1 AND payment_status = 'pending';

-- name: FailPayment :exec
UPDATE payments
  set payment_status = 'failed', payment_date = NOW()
WHERE id = $1 AND payment_status = 'pending';

-- name: GetExpiredPayments :many
SELECT reservation_id FROM payments
WHERE payment_status = 'pending' AND creaate_at < NOW() - INTERVAL '15 minutes';

-- name: DeletePayment :exec
DELETE FROM payments
WHERE id = $1;
