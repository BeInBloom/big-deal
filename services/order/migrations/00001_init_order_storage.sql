-- +goose Up

CREATE TYPE order_status AS ENUM (
    'PENDING_PAYMENT',
    'PAID',
    'CANCELLED'
);

CREATE TYPE payment_method AS ENUM (
    'CARD',
    'SBP',
    'CREDIT_CARD',
    'INVESTOR_MONEY'
);

CREATE TABLE orders (
    id uuid PRIMARY KEY,
    user_id uuid NOT NULL,
    status order_status NOT NULL,
    payment_method payment_method,
    transaction_id uuid,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT orders_payment_state_check CHECK (
        (
            status = 'PAID'
            AND payment_method IS NOT NULL
            AND transaction_id IS NOT NULL
        )
        OR
        (
            status IN ('PENDING_PAYMENT', 'CANCELLED')
            AND payment_method IS NULL
            AND transaction_id IS NULL
        )
    )
);

CREATE TABLE order_parts (
    order_id uuid NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    part_id uuid NOT NULL,
    description text NOT NULL,
    price bigint NOT NULL CHECK (price >= 0),

    PRIMARY KEY (order_id, part_id)
);

CREATE INDEX orders_user_id_idx ON orders(user_id);
CREATE INDEX order_parts_part_id_idx ON order_parts(part_id);
