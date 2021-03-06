CREATE TABLE managers (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    salary INTEGER NOT NULL CHECK (salary > 0),
    plan INTEGER NOT NULL CHECK (plan > 0),
    boss_id BIGINT NOT NULL,
    department TEXT,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    login TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
)