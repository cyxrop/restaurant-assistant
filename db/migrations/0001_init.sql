-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE EXTENSION IF NOT EXISTS pgcrypto;

DROP TYPE IF EXISTS user_account_type;
CREATE TYPE user_account_type AS ENUM ('default', 'admin');

CREATE TABLE IF NOT EXISTS user_account
(
    id          UUID         NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    name        VARCHAR(128) UNIQUE                    NOT NULL,
    password    VARCHAR(256)                           NOT NULL,
    email       VARCHAR(128)                           NOT NULL,
    description VARCHAR(512)                           NOT NULL,
    user_type   user_account_type DEFAULT 'default'    NOT NULL,
    update_time TIMESTAMPTZ  DEFAULT NOW()             NOT NULL,
    create_time TIMESTAMPTZ  DEFAULT CLOCK_TIMESTAMP() NOT NULL
);

CREATE SEQUENCE IF NOT EXISTS product_id_seq;

CREATE TABLE IF NOT EXISTS product
(
    id          INTEGER NOT NULL DEFAULT nextval('product_id_seq') PRIMARY KEY,
    name        VARCHAR(128) UNIQUE                    NOT NULL,
    description VARCHAR(512)                           NOT NULL,
    update_time TIMESTAMPTZ  DEFAULT NOW()             NOT NULL,
    create_time TIMESTAMPTZ  DEFAULT CLOCK_TIMESTAMP() NOT NULL
);

CREATE SEQUENCE IF NOT EXISTS customer_order_id_seq;

CREATE TABLE IF NOT EXISTS customer_order
(
    id          INTEGER NOT NULL DEFAULT nextval('customer_order_id_seq') PRIMARY KEY,
    order_data  jsonb NOT NULL,
    user_id     UUID REFERENCES user_account (id) ON DELETE CASCADE NOT NULL,
    update_time TIMESTAMPTZ  DEFAULT NOW()             NOT NULL,
    create_time TIMESTAMPTZ  DEFAULT CLOCK_TIMESTAMP() NOT NULL
);

-- Fill user_account table
INSERT INTO user_account VALUES ('6ba7b812-9dad-11d1-80b4-00c04fd430c8', 'test', '$2a$14$6MIB81kfk8wJTpL169510ObxqflRrcAMltL/awXXqXLHu.1gg24iS', 'email', 'desc', 'admin');

-- Fill product table
INSERT INTO product VALUES (DEFAULT , 'cake', 'desc'), (DEFAULT , 'tea', 'desc');

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TYPE IF EXISTS user_account_type CASCADE;

DROP TABLE IF EXISTS customer_order;

DROP TABLE IF EXISTS product;

DROP TABLE IF EXISTS user_account;

DROP SEQUENCE IF EXISTS product_id_seq;

DROP SEQUENCE IF EXISTS customer_order_id_seq;
