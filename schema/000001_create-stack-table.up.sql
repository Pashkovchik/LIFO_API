BEGIN;
CREATE TABLE IF NOT EXISTS "stack" (
    id           BIGSERIAL    PRIMARY KEY,
    name         VARCHAR,
    created_date   TIMESTAMP   NOT NULL DEFAULT NOW()
    );
COMMIT;