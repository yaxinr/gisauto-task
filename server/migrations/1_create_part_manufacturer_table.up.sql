CREATE TABLE part_manufacturer (
    id SERIAL PRIMARY KEY,
    name varchar(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
)