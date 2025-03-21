CREATE schema 'events'

CREATE TABLE events.events (
                        id UUID PRIMARY KEY,
                        title TEXT NOT NULL,
                        start_time TIMESTAMP NOT NULL,
                        duration INTERVAL NOT NULL,
                        description TEXT,
                        user_id UUID NOT NULL,
                        notify_before INTERVAL
);

CREATE INDEX idx_events_start_time ON events(start_time);
