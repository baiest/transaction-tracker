ALTER TABLE movements ADD COLUMN new_id VARCHAR(255) NOT NULL;

ALTER TABLE movements ADD CONSTRAINT movements_new_id_key UNIQUE (new_id);

ALTER TABLE movements DROP CONSTRAINT movements_pkey;
ALTER TABLE movements ADD CONSTRAINT movements_pkey PRIMARY KEY (new_id);

ALTER TABLE movements DROP COLUMN id;

ALTER TABLE movements RENAME COLUMN new_id TO id;
