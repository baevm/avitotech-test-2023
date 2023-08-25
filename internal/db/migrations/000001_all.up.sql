CREATE TABLE IF NOT EXISTS segments (
    id bigserial PRIMARY KEY,
    slug varchar (255) NOT NULL UNIQUE,
    created_at timestamp NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS user_segments (
    segment_id bigserial references segments(id) NOT NULL,
    user_id bigint NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),

    PRIMARY KEY (user_id, segment_id)
);