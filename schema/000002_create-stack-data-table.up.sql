BEGIN;
CREATE TABLE IF NOT EXISTS "stack_data" (
    id           BIGSERIAL    PRIMARY KEY,
    stack_id     INT,
    info         VARCHAR,
    created_date   TIMESTAMP   NOT NULL DEFAULT NOW(),


    CONSTRAINT fk_stack FOREIGN KEY (stack_id) REFERENCES stack (id)
    );
COMMIT;