-- +goose Up

CREATE TABLE feed_follows(
     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID   NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    feed_id UUID    NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    CONSTRAINT unique_user_feed_follow UNIQUE (user_id, feed_id)


);

-- +goose Down

DROP TABLE feed_follows;