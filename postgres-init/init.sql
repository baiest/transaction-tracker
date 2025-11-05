-- Ensure database and extension exist
CREATE DATABASE transactions;

\connect transactions;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
