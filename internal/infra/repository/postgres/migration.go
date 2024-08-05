package postgres

// TODO переделать на нормальные миграции
// https://github.com/pressly/goose
var migrationQuery = `
CREATE TYPE user_role AS ENUM ('client', 'moderator');

CREATE TABLE IF NOT EXISTS user (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email VARCHAR(255) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  type user_role NOT NULL
);

CREATE TABLE subscription (
  house_id BIGINT NOT NULL,
  email VARCHAR(255) NOT NULL,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
  PRIMARY KEY (house_id, email),
  CONSTRAINT fk_subscription_house_id FOREIGN KEY (house_id) REFERENCES house (id) ON DELETE CASCADE,
  CONSTRAINT fk_subscription_email FOREIGN KEY (email) REFERENCES user (email) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS house (
  id BIGSERIAL PRIMARY KEY,
  address VARCHAR(255) NOT NULL,
  year_of_construction BIGINT NOT NULL,
  developer VARCHAR(255),
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TYPE moderation_status AS ENUM ('created', 'on_moderation', 'approved', 'declined');

CREATE TABLE IF NOT EXISTS flat (
  id BIGSERIAL PRIMARY KEY,
  house_id BIGINT NOT NULL,
  price BIGINT NOT NULL CHECK (price > 0),
  rooms BIGINT NOT NULL CHECK (rooms > 0),
  status moderation_status NOT NULL DEFAULT 'created',
  CONSTRAINT fk_flat_house_id FOREIGN KEY (house_id) REFERENCES house (id) ON DELETE CASCADE
);
`
