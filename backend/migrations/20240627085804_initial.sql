-- +goose Up
CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    username TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    admin BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE users_telegrams (
    id uuid REFERENCES users (id),
    chat_id BIGINT NOT NULL UNIQUE,
    telegram_id BIGINT NOT NULL UNIQUE
);

CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

INSERT INTO
    categories (name)
VALUES ('Горячее'),
    ('Холодное'),
    ('Напиток'),
    ('Острое'),
    ('Рыба'),
    ('Вегетарианское'),
    ('Мясное');

CREATE TABLE dish (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    -- Цена в минимальных единицах валюты, например, в копейках
    price INT NOT NULL CHECK (price > 0),
    image_id TEXT
);

CREATE TABLE dish_categories (
    dish_id INT REFERENCES dish (id),
    category_id INT REFERENCES categories (id),
    PRIMARY KEY (dish_id, category_id)
);

CREATE TABLE orders (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    payment_method TEXT NOT NULL,
    user_id uuid NOT NULL REFERENCES users (id),
    total BIGINT NOT NULL CHECK (total > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
    wishes TEXT
);

CREATE TABLE order_items (
    order_id uuid NOT NULL REFERENCES orders (id),
    dish_id INT NOT NULL REFERENCES dish (id),
    status TEXT NOT NULL,
    count INT NOT NULL CHECK (count > 0),
    price INT NOT NULL CHECK (price > 0),
    PRIMARY KEY (order_id, dish_id)
);

-- +goose Down
DROP TABLE order_items;

DROP TABLE orders;

DROP TABLE dish_categories;

DROP TABLE dish;

DROP TABLE categories;

DROP TABLE users_telegrams;

DROP TABLE users;