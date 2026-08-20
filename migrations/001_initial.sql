CREATE TABLE sites (id TEXT PRIMARY KEY, name TEXT NOT NULL, location TEXT, created_at TIMESTAMPTZ NOT NULL, version BIGINT NOT NULL);
CREATE TABLE devices (id TEXT PRIMARY KEY, site_id TEXT NOT NULL REFERENCES sites(id), name TEXT NOT NULL, protocol TEXT NOT NULL, address TEXT NOT NULL, enabled BOOLEAN NOT NULL, version BIGINT NOT NULL);
CREATE TABLE points (id TEXT PRIMARY KEY, device_id TEXT NOT NULL REFERENCES devices(id), name TEXT NOT NULL, address INTEGER NOT NULL, data_type TEXT NOT NULL, unit TEXT, scale DOUBLE PRECISION NOT NULL, deadband DOUBLE PRECISION NOT NULL, sample_every INTERVAL NOT NULL, version BIGINT NOT NULL);
CREATE TABLE telemetry (point_id TEXT NOT NULL, observed_at TIMESTAMPTZ NOT NULL, value DOUBLE PRECISION NOT NULL, quality TEXT NOT NULL, sequence BIGINT NOT NULL, PRIMARY KEY(point_id, observed_at, sequence));
CREATE INDEX telemetry_point_time ON telemetry(point_id, observed_at DESC);
