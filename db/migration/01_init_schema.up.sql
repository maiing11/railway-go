CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
  id uuid PRIMARY KEY,
  name VARCHAR NOT NULL,
  email VARCHAR UNIQUE NOT NUll,
  password VARCHAR NOT NUll,
  phoneNumber VARCHAR NOT NUll,
  created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE passengers (
  id bigserial PRIMARY KEY,
  name VARCHAR NOT NULL,
  id_number VARCHAR NOT NULL,
  user_id uuid DEFAULT NULL,
  created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE trains (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR NOT NULL,
  capacity int not NULL,
  created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE routes (
  id BIGSERIAL PRIMARY KEY,
  source_station text NOT NULL,
  destination_station text NOT NULL,
  travel_time int NOT NULL,
  created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE tipe_class AS ENUM ('premium', 'economy', 'luxury');

CREATE TABLE schedules (
  id BIGSERIAL PRIMARY KEY,
  train_id int NOT NUll,
  class_type tipe_class NOT NUll,
  departure_date timestamp NOT NUll,
  arrival_date timestamp NOT NUll,
  available_seats int NOT NUll,
  price BIGINT NOT NUll,
  route_id int NOT NUll,
  created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (train_id) REFERENCES trains(id),
  FOREIGN KEY (route_id) REFERENCES routes(id)
);

CREATE TYPE status_reservation AS ENUM ('pending', 'canceled', 'success');

CREATE TABLE reservations (
  id uuid PRIMARY KEY,
  passenger_id bigint NOT NUll,
  schedule_id bigint NOT NUll,
  seat_number INTEGER,
  booking_date timestamp NOT NUll,
  payment_id uuid NOT NUll,
  status status_reservation DEFAULT 'pending',
  created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (passenger_id) REFERENCES passengers(id),
  FOREIGN KEY (schedule_id) REFERENCES schedules(id)
);

CREATE TABLE payments (
  id uuid PRIMARY KEY,
  reservation_id uuid,
  payment_method text,
  amount bigint,
  transaction_id text,
  payment_date timestamp,
  gateway_response text,
  status text,
  created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (reservation_id) REFERENCES reservations(id)
);
