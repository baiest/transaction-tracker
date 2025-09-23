ALTER TABLE movements ADD COLUMN old_id UUID DEFAULT gen_random_uuid();

UPDATE movements SET old_id = gen_random_uuid();

ALTER TABLE movements DROP CONSTRAINT movements_pkey;
ALTER TABLE movements ADD CONSTRAINT movements_pkey PRIMARY KEY (old_id);

ALTER TABLE movements DROP COLUMN id;

ALTER TABLE movements RENAME COLUMN old_id TO id;
