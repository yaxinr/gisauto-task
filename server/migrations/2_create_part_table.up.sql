CREATE TABLE part (
    id SERIAL PRIMARY KEY,
    manufacturer_id integer,
    vendor_code varchar(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
)