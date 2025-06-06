CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE user_role AS ENUM ('user', 'general affairs', 'admin');

-- 🧑 Users Table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR NOT NULL,
    email VARCHAR UNIQUE NOT NULL,
    password VARCHAR NOT NULL,
    phone_number VARCHAR NOT NULL,
    role user_role DEFAULT 'user',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- 🎟️ Passengers Table (Linked to Users)
CREATE TABLE passengers (
    id UUID  PRIMARY KEY,
    name VARCHAR NOT NULL,
    id_number VARCHAR UNIQUE NOT NULL,
    user_id UUID,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_passenger_user
ON passengers (user_id);

-- 🚆 Trains Table
CREATE TABLE trains (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    capacity INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stations (
  id BIGSERIAL PRIMARY KEY,
  code VARCHAR (4) NOT NULL UNIQUE,
  station_name VARCHAR NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- 🏁 Routes Table
CREATE TABLE routes (
    id BIGSERIAL PRIMARY KEY,
    source_station VARCHAR (4) NOT NULL,
    destination_station VARCHAR(4) NOT NULL,
    travel_time INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (source_station) REFERENCES stations(code) ON DELETE CASCADE,
    FOREIGN KEY (destination_station) REFERENCES stations(code) ON DELETE CASCADE
);

-- 🚉 Train Schedules Table
CREATE TYPE tipe_class AS ENUM ('premium', 'economy', 'luxury');

CREATE TABLE wagons (
  id BIGSERIAL PRIMARY KEY,
  train_id BIGINT NOT NULL,
  wagon_number INT NOT NULL CHECK (wagon_number > 0),
  class_type tipe_class NOT NULL,
  total_seats INT NOT NULL CHECK (total_seats <= 25),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (train_id) REFERENCES trains(id) ON DELETE CASCADE
);

CREATE TABLE discount_codes (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  code TEXT UNIQUE NOT NULL,
  discount_percent INT NOT NULL CHECK (discount_percent BETWEEN 1 AND 100),
  max_uses INT NOT NULL CHECK (max_uses >= 0),
  expires_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE seat_row AS ENUM ('A', 'B', 'C', 'D');

CREATE TABLE seats (
  id BIGSERIAL PRIMARY KEY,
  wagon_id BIGINT,
  seat_number INT NOT NULL CHECK (seat_number BETWEEN 1 AND 25),
  seat_row seat_row NOT NULL,
  is_available BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (wagon_id) REFERENCES wagons(id) ON DELETE CASCADE
);


CREATE TABLE schedules (
  id BIGSERIAL PRIMARY KEY,
  train_id BIGINT NOT NULL,
  route_id BIGINT NOT NULL,
  departure_date TIMESTAMP NOT NULL,
  arrival_date TIMESTAMP NOT NULL,
  available_seats INT NOT NULL,
  price BIGINT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (train_id) REFERENCES trains(id) ON DELETE CASCADE,
  FOREIGN KEY (route_id) REFERENCES routes(id) ON DELETE CASCADE
);

CREATE TYPE status_reservation AS ENUM ('pending', 'success', 'cancelled');

CREATE TABLE reservations (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  passenger_id UUID NOT NULL,
  schedule_id BIGINT NOT NULL,
  wagon_id BIGINT NOT NULL,
  seat_id BIGINT NOT NULL,
  booking_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  discount_id UUID,
  price BIGINT,
  reservation_status status_reservation NOT NULL DEFAULT 'pending',
  expires_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '15 minutes',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (passenger_id) REFERENCES passengers(id) ON DELETE CASCADE,
  FOREIGN KEY (schedule_id) REFERENCES schedules(id) ON DELETE CASCADE,
  FOREIGN KEY (wagon_id) REFERENCES wagons(id) ON DELETE CASCADE,
  FOREIGN KEY (seat_id) REFERENCES seats(id) ON DELETE CASCADE,
  FOREIGN KEY (discount_id) REFERENCES discount_codes(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX unique_active_seat_reservation
ON reservations (schedule_id, wagon_id, seat_id)
WHERE reservation_status IN ('pending', 'success');

CREATE TABLE payments (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  reservation_id UUID NOT NULL,
  payment_method TEXT NOT NULL,
  amount BIGINT NOT NULL,
  transaction_id TEXT UNIQUE NOT NULL,
  payment_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  gateway_response TEXT,
  payment_status TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (reservation_id) REFERENCES reservations(id) ON DELETE CASCADE
);



-- Applied Discounts (track discount used in reservations)
CREATE TABLE reservation_discounts (
  reservation_id UUID not null,
  discount_id UUID not null,
  primary key (reservation_id, discount_id),
  foreign key (reservation_id) references reservations(id) on delete cascade,
  foreign key (discount_id) references discount_codes(id) on delete cascade
)

