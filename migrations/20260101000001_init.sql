-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users
(
    id            BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    login         VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    CONSTRAINT uq_users_login UNIQUE (login)
);

CREATE TABLE IF NOT EXISTS orders
(
    id          BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id     BIGINT       NOT NULL REFERENCES users (id),
    number      VARCHAR(255) NOT NULL,
    status      VARCHAR(50)  NOT NULL DEFAULT 'NEW',
    accrual     DOUBLE PRECISION,
    uploaded_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    CONSTRAINT uq_orders_number UNIQUE (number)
);

CREATE TABLE IF NOT EXISTS withdrawals
(
    id           BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id      BIGINT           NOT NULL REFERENCES users (id),
    order_number VARCHAR(255)     NOT NULL,
    sum          DOUBLE PRECISION NOT NULL,
    processed_at TIMESTAMPTZ      NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS withdrawals;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
