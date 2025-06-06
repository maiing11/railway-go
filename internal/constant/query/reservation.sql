-- name: GetReservation :one
SELECT * FROM reservations
WHERE id = $1 LIMIT 1;

-- name: ListReservations :many
SELECT 
r.id AS reservation_id,
p.name AS passenger_name,
p.id_number AS passenger_id_number,
u.name AS user_name,
u.email AS user_email,
s.departure_date,
s.arrival_date,
r.price AS ticket_price,
t.name AS train_name,
w.class_type,
w.wagon_number,
st.seat_number,
st.seat_row,
r.booking_date,
r.reservation_status,
rt.source_station,
rt.destination_station,
d.code AS discount_code,
d.discount_percent,
py.amount AS payment_amount,
py.payment_method,
py.payment_status
FROM reservations r
LEFT JOIN passengers p ON r.passenger_id = p.id
LEFT JOIN users u ON p.user_id = u.id
LEFT JOIN schedules s ON r.schedule_id = s.id
LEFT JOIN seats st ON r.seat_id = st.id
LEFT JOIN wagons w ON r.wagon_id = w.id
LEFT JOIN trains t ON s.train_id = t.id
LEFT JOIN routes rt ON s.route_id = rt.id
LEFT JOIN discount_codes d ON r.discount_id = d.id
LEFT JOIN payments py ON r.id = py.reservation_id
ORDER BY booking_date DESC
LIMIT $1
OFFSET $2;

-- name: CountReservations :one
select count(*) from reservations;

-- name: CreateReservation :one
INSERT INTO reservations (
   passenger_id, schedule_id, wagon_id, seat_id, booking_date, reservation_status, discount_id, price, expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) 
RETURNING *;

-- name: UpdateReservation :exec
UPDATE reservations
  set  passenger_id = $2 , schedule_id = $3, wagon_id=$4, seat_id = $5, booking_date = $6, reservation_status = $7, discount_id = $8, price = $9, expires_at = $10, updated_at = NOW()
WHERE id = $1;


-- name: DeleteReservation :exec
DELETE FROM reservations
WHERE id = $1 AND reservation_status = 'canceled';

-- name: CheckSeatAvailability :one
SELECT COUNT(*) FROM reservations 
WHERE schedule_id = $1 AND wagon_id = $2 AND seat_id = $3
  AND reservation_status IN ('pending', 'success');


-- -- name: HoldSeat :exec
-- INSERT INTO seat_holds (passenger_id, schedule_id, wagon_id, seat_id, expires_at)
-- VALUES ($1, $2, $3, $4, NOW() + INTERVAL '15 minutes')
-- RETURNING *;

-- -- name: CreateReservationFromHold :one
-- WITH deleted_hold AS (
--     DELETE FROM seat_holds
--     WHERE seat_holds.passenger_id = $1 
--     AND seat_holds.schedule_id = $2
--     AND seat_holds.wagon_id = $3
--     AND seat_holds.seat_id = $4
--     RETURNING passenger_id, schedule_id, wagon_id, seat_id
-- )
-- INSERT INTO reservations (
--     id, passenger_id, schedule_id, wagon_id, seat_id, booking_date, reservation_status, discount_id, price, expires_at
-- -- )
-- SELECT 
--     uuid_generate_v4(), deleted_hold.passenger_id, deleted_hold.schedule_id, deleted_hold.wagon_id, deleted_hold.seat_id,
--     NOW(), 'pending', $5, $6, NOW() + INTERVAL '15 minutes'
-- FROM deleted_hold
-- RETURNING *;

-- name: ConfirmReservation :exec
UPDATE reservations
SET reservation_status = 'success', updated_at = NOW()
WHERE id = $1 AND reservation_status = 'pending';

-- name: CancelReservation :exec
UPDATE reservations
SET reservation_status = 'cancelled', updated_at = NOW()
WHERE id = $1 AND reservation_status = 'pending';


-- name: UpdateTrainCapacity :exec
UPDATE schedules
SET available_seats = available_seats - 1
WHERE id = $1 AND available_seats > 0;

-- -- name: CleanupExpiredHolds :exec
-- DELETE FROM seat_holds WHERE expires_at < NOW();

-- -- name: ExpireSeatHolds :exec
-- DELETE FROM seat_holds
-- WHERE expires_at < NOW();

-- name: ExpireUndpaidReservations :exec
DELETE FROM reservations
WHERE expires_at < NOW() AND reservation_status = 'pending';


-- name: GetFullReservation :one
SELECT 
r.id AS reservation_id,
p.name AS passenger_name,
p.id_number AS passenger_id_number,
u.name AS user_name,
u.email AS user_email,
s.departure_date,
s.arrival_date,
r.price AS ticket_price,
t.name AS train_name,
w.class_type,
w.wagon_number,
st.seat_number,
st.seat_row,
r.booking_date,
r.reservation_status,
rt.source_station,
rt.destination_station,
d.code AS discount_code,
d.discount_percent,
py.amount AS payment_amount,
py.payment_method,
py.payment_status 
FROM reservations r
LEFT JOIN passengers p ON r.passenger_id = p.id
LEFT JOIN users u ON p.user_id = u.id
LEFT JOIN schedules s ON r.schedule_id = s.id
LEFT JOIN seats st ON r.seat_id = st.id
LEFT JOIN wagons w ON r.wagon_id = w.id
LEFT JOIN trains t ON s.train_id = t.id
LEFT JOIN routes rt ON s.route_id = rt.id
LEFT JOIN discount_codes d ON r.discount_id = d.id
LEFT JOIN payments py ON r.id = py.reservation_id
WHERE r.id = $1;
