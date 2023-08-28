CREATE TABLE IF NOT EXISTS segments (
    slug varchar (255) PRIMARY KEY,
    user_percent int,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS user_segments (
    segment_slug varchar (255) references segments(slug) ON DELETE CASCADE,
    user_id bigserial NOT NULL references users(id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT now(),
    expire_at timestamptz,

    PRIMARY KEY (user_id, segment_slug)
);

CREATE TABLE IF NOT EXISTS user_segment_history (
    id bigserial PRIMARY KEY,
    segment_slug varchar (255) NOT NULL,
    user_id bigserial NOT NULL,
    operation varchar(1) NOT NULL, -- I for insert, D for delete
    executed_at timestamptz NOT NULL DEFAULT now()
);

CREATE OR REPLACE FUNCTION user_segments_trigger()
RETURNS TRIGGER AS $$
BEGIN
    -- Insert a row into user_segment_history when a row is inserted or deleted
    IF TG_OP = 'INSERT' THEN
        INSERT INTO user_segment_history (segment_slug, user_id, operation, executed_at)
        VALUES (NEW.segment_slug, NEW.user_id, 'I', now());
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        INSERT INTO user_segment_history (segment_slug, user_id, operation, executed_at)
        VALUES (OLD.segment_slug, OLD.user_id, 'D', now());
        RETURN OLD;
    END IF;
    RETURN NULL; -- Return NULL for other operations
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER user_segments_history_trigger
AFTER INSERT OR DELETE ON user_segments
FOR EACH ROW
EXECUTE FUNCTION user_segments_trigger();