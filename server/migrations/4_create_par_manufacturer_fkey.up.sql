ALTER TABLE part
    ADD CONSTRAINT manufacturer_fkey FOREIGN KEY (manufacturer_id)
    REFERENCES part_manufacturer (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID
