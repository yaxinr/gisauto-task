CREATE INDEX part_manufacturer_name
    ON part_manufacturer USING hash
    (name)
;