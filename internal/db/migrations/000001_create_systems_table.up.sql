CREATE TYPE status AS ENUM ('operational', 'degraded', 'down');

CREATE TABLE IF NOT EXISTS systems (
    id      SERIAL PRIMARY KEY,
    name    VARCHAR(255) NOT NULL UNIQUE,
    current_status status NOT NULL DEFAULT 'operational'
);

CREATE TABLE IF NOT EXISTS incidents (
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    details     TEXT,
    created_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS system_incidents (
    id            SERIAL PRIMARY KEY,
    system_id     INTEGER REFERENCES systems(id) ON DELETE CASCADE,
    incident_id   INTEGER REFERENCES incidents(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS updates (
    id              SERIAL PRIMARY KEY,
    incident_id     INTEGER REFERENCES incidents(id),
    title           VARCHAR(255) NOT NULL,
    details         TEXT,
    created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    new_status      status NOT NULL DEFAULT 'operational'
);